package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"quillcrypt-backend/internal/config"
	"quillcrypt-backend/internal/core/domain"
	"quillcrypt-backend/internal/core/port"
	"quillcrypt-backend/internal/repository/redis"
	"quillcrypt-backend/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	userService port.UserService
}

func NewAuthHandler(userService port.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) BeginAuth(c fiber.Ctx) error {
	state := c.Query("state")
	challenge := c.Query("code_challenge")

	if state == "" || challenge == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": "state and code_challenge are required"})
	}

	key := fmt.Sprintf("auth_state:%s", state)
	err := redis.Client.Set(c.Context(), key, challenge, 10*time.Minute).Err()
	if err != nil {
		logger.Error("Redis set error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": fiber.StatusInternalServerError, "message": http.StatusText(fiber.StatusInternalServerError)})
	}

	url := config.OAuth2Config.AuthCodeURL(state,

		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	return c.Redirect().To(url)
}

func (h *AuthHandler) AuthCallback(c fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": "code and state are required"})
	}

	key := fmt.Sprintf("auth_state:%s", state)
	exists, err := redis.Client.Exists(c.Context(), key).Result()
	if err != nil || exists == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": fiber.StatusForbidden, "message": "Invalid or expired state"})
	}

	redirectURL, err := url.Parse(config.Config.MobileCallback)
	if err != nil {
		logger.Error("Cannot parse url from QC_MOBILECALLBACK")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": http.StatusText(fiber.StatusInternalServerError),
		})
	}
	params := url.Values{}
	params.Add("state", state)
	params.Add("code", code)
	redirectURL.RawQuery = params.Encode()
	return c.Redirect().To(redirectURL.String())
}

func (h *AuthHandler) ExchangeAuth(c fiber.Ctx) error {
	var req struct {
		Code         string `json:"code"`
		State        string `json:"state"`
		CodeVerifier string `json:"code_verifier"`
	}

	if err := c.Bind().Body(&req); err != nil {
		logger.Error("Exchange error", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": http.StatusText(fiber.StatusBadRequest)})
	}

	key := fmt.Sprintf("auth_state:%s", req.State)
	storedChallenge, err := redis.Client.Get(c.Context(), key).Result()
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": fiber.StatusForbidden, "message": "Expired or invalid state"})
	}

	hash := sha256.Sum256([]byte(req.CodeVerifier))
	localChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	if localChallenge != storedChallenge {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": "Fraudulent verifier"})
	}

	token, err := config.OAuth2Config.Exchange(c.Context(), req.Code, oauth2.SetAuthURLParam("code_verifier", req.CodeVerifier))
	if err != nil {
		logger.Error("Token exchange error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": fiber.StatusInternalServerError, "message": http.StatusText(fiber.StatusInternalServerError)})
	}

	redis.Client.Del(c.Context(), key)

	client := config.OAuth2Config.Client(c.Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		logger.Error("Get GitHub user error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": fiber.StatusInternalServerError, "message": http.StatusText(fiber.StatusInternalServerError)})
	}
	defer resp.Body.Close()

	var ghUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": fiber.StatusInternalServerError, "message": http.StatusText(fiber.StatusInternalServerError)})
	}

	if ghUser.Email == "" {
		ghUser.Email = h.fetchGithubEmail(client)
	}

	user := &domain.User{
		ID:        ghUser.ID,
		Username:  ghUser.Login,
		Email:     ghUser.Email,
		AvatarURL: ghUser.AvatarURL,
	}
	if user.Username == "" {
		user.Username = ghUser.Name
	}

	dbUser, err := h.userService.RegisterOrLogin(c.Context(), user)
	if err != nil {
		logger.Error("Service RegisterOrLogin error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": http.StatusText(fiber.StatusInternalServerError),
		})
	}

	sess, err := redis.Store.Get(c)
	if err != nil {
		logger.Error("Session error", zap.Error(err))
	} else {
		sess.Set("user_id", dbUser.ID)
		if err := sess.Save(); err != nil {
			logger.Error("Session save error", zap.Error(err))
		}
	}

	return c.Status(fiber.StatusOK).JSON(dbUser)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	sess, err := redis.Store.Get(c)
	if err == nil {
		sess.Destroy()
	}

	c.Cookie(&fiber.Cookie{
		Name:     "qc_session",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": fiber.StatusOK, "message": "Logged out"})
}

func (h *AuthHandler) fetchGithubEmail(client *http.Client) string {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return ""
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}
	if len(emails) > 0 {
		return emails[0].Email
	}
	return ""
}
