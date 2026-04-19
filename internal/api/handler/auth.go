package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"quillcrypt-backend/internal/config"
	"quillcrypt-backend/internal/repository/redis"
	"quillcrypt-backend/pkg/logger"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

func BeginAuth(c fiber.Ctx) error {
	provider := c.Params("provider")
	logger.Debug("cb", zap.String("cb_url", config.Config.Google_Callback))
	handler := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// gothic looks for "provider" in query or session
		q := r.URL.Query()
		q.Set("provider", provider)
		r.URL.RawQuery = q.Encode()

		if user, err := gothic.CompleteUserAuth(w, r); err == nil {
			logger.Debug("Logged in user", zap.Any("User", user))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})
	return handler(c)
}

func AuthCallback(c fiber.Ctx) error {
	provider := c.Params("provider")
	
	handler := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Set("provider", provider)
		r.URL.RawQuery = q.Encode()

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			logger.Error("Auth Callback error", zap.Error(err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(map[string]any{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		
		logger.Debug("Logged in user", zap.Any("User", user))

		// Save to Redis Session
		sess, err := redis.Store.Get(c)
		if err == nil {
			sess.Set("user_id", user.UserID)
			if err := sess.Save(); err != nil {
				logger.Error("Session save error", zap.Error(err))
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return handler(c)
}

func Logout(c fiber.Ctx) error {
	sess, err := redis.Store.Get(c)
	if err == nil {
		if err := sess.Destroy(); err != nil {
			logger.Error("Fiber session destroy error", zap.Error(err))
		}
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
