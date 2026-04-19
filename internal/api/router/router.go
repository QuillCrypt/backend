package router

import (
	"quillcrypt-backend/internal/api/handler"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Get("/:provider", handler.BeginAuth)
	auth.Get("/:provider/callback", handler.AuthCallback)
	auth.Post("/logout", handler.Logout)
}
