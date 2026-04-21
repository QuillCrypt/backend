package handler

import (
	"encoding/json"
	"net/http"
	"quillcrypt-backend/internal/config"
	"quillcrypt-backend/internal/core/domain"
	"quillcrypt-backend/internal/core/port"
	"quillcrypt-backend/internal/repository/redis"
	"quillcrypt-backend/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userService port.UserService
}

func NewAuthHandler(userService port.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) BeginAuth(c fiber.Ctx) error {
	provider := c.Params("provider")
	logger.Debug("cb", zap.String("cb_url", config.Config.Google_Callback))
	handler := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Set("provider", provider)
		r.URL.RawQuery = q.Encode()

		if _, err := gothic.CompleteUserAuth(w, r); err == nil {
			w.WriteHeader(fiber.StatusOK)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})
	return handler(c)
}

func (h *AuthHandler) AuthCallback(c fiber.Ctx) error {
	provider := c.Params("provider")

	handler := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Set("provider", provider)
		r.URL.RawQuery = q.Encode()

		gothUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			logger.Error("Auth Callback error", zap.Error(err))
			w.WriteHeader(fiber.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"error": "Auth failed"})
			return
		}

		user := &domain.User{
			Username:  gothUser.NickName,
			Email:     gothUser.Email,
			AvatarURL: gothUser.AvatarURL,
		}
		switch provider {
		case "google":
			user.GoogleID = &gothUser.UserID
		case "github":
			user.GithubID = &gothUser.UserID
		}
		if user.Username == "" {
			user.Username = gothUser.Name
		}
		dbUser, err := h.userService.RegisterOrLogin(r.Context(), user)
		if err != nil {
			logger.Error("Service RegisterOrLogin error", zap.Error(err))
			w.WriteHeader(fiber.StatusInternalServerError)
			return
		}

		sess, err := redis.Store.Get(c)
		if err == nil {
			sess.Set("user_id", dbUser.ID.String())
			if err := sess.Save(); err != nil {
				logger.Error("Session save error", zap.Error(err))
			}
		}

		w.WriteHeader(fiber.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dbUser)
	})
	return handler(c)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	sess, err := redis.Store.Get(c)
	if err == nil {
		sess.Destroy()
	}

	handler := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
	})

	c.Cookie(&fiber.Cookie{
		Name:     "qc_session",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})

	return handler(c)
}
