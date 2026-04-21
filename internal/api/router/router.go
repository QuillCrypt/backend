package router

import (
	"quillcrypt-backend/internal/api/handler"
	"quillcrypt-backend/internal/api/middleware"
	"quillcrypt-backend/internal/core/port"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, userService port.UserService) {
	app.Use(middleware.WithLogs())

	authHandler := handler.NewAuthHandler(userService)

	auth := app.Group("/auth")
	auth.Get("/:provider", authHandler.BeginAuth)
	auth.Get("/:provider/callback", authHandler.AuthCallback)
	auth.Post("/logout", authHandler.Logout)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Route not found",
		})
	})
}
