package middleware

import (
	"github.com/gofiber/fiber/v3"
	"quillcrypt-backend/internal/repository/redis"
)

func WithAuth(c fiber.Ctx) error {
	sess, err := redis.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Unauthorized: Session error",
		})
	}

	userId := sess.Get("user_id")
	if userId == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Unauthorized: Login required",
		})
	}

	c.Locals("user_id", userId)

	return c.Next()
}
