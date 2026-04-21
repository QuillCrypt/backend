package middleware

import (
	"os"
	"quillcrypt-backend/internal/config"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func WithLogs() fiber.Handler {
	return logger.New(logger.Config{
		Done: func(c fiber.Ctx, log []byte) {
			if config.Config.LogFilePath != "" {
				os.WriteFile(config.Config.LogFilePath, log, os.ModeAppend)
				os.WriteFile(config.Config.LogFilePath, []byte("\n"), os.ModeAppend)
			}
		},
	})
}
