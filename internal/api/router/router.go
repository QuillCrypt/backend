package router

import (
	"net/http"
	"quillcrypt-backend/internal/api/handler"
	"quillcrypt-backend/internal/api/middleware"
	"quillcrypt-backend/internal/core/port"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, userService port.UserService) {
	app.Use(middleware.WithLogs())

	authHandler := handler.NewAuthHandler(userService)

	auth := app.Group("/auth")
	auth.Get("/", authHandler.BeginAuth)
	auth.Get("/callback", authHandler.AuthCallback)
	auth.Post("/exchange", authHandler.ExchangeAuth)
	auth.Post("/logout", authHandler.Logout)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": http.StatusText(fiber.StatusNotFound),
		})
	})
}
