package middleware

import (
	"net/http"
	"quillcrypt-backend/internal/repository/redis"
	"quillcrypt-backend/pkg/logger"

	"github.com/gofiber/fiber/v3"
)

func WithAuth(c fiber.Ctx) error {
	sess, err := redis.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": http.StatusText(fiber.StatusUnauthorized),
		})
	}

	userID, ok := sess.Get("user_id").(int64)
	if !ok {
		logger.Error("cannot convert user_id from session storage to int64")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": http.StatusText(fiber.StatusInternalServerError),
		})
	}
	if userID <= 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": http.StatusText(fiber.StatusUnauthorized),
		})
	}

	c.Locals("user_id", userID)

	return c.Next()
}
