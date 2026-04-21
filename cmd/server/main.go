package main

import (
	"net"
	"os"
	"os/signal"
	"quillcrypt-backend/internal/api/router"
	"quillcrypt-backend/internal/config"
	"quillcrypt-backend/internal/core/service"
	"quillcrypt-backend/internal/repository/postgres"
	"quillcrypt-backend/internal/repository/redis"
	"quillcrypt-backend/pkg/logger"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func init() {
	config.LoadConfig()
}

func main() {
	logger.Init(config.Config.Mode)
	defer logger.Log.Sync()

	redis.InitSession()
	defer redis.Store.Storage.Close()

	postgres.InitDB()
	defer postgres.DB.Close()
	app := fiber.New()

	userRepo := postgres.NewUserRepository(postgres.DB)
	userService := service.NewUserService(userRepo)
	router.SetupRoutes(app, userService)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		if !fiber.IsChild() {
			logger.Info("Shutting down server...")
		}
		if err := app.Shutdown(); err != nil {
			logger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	addr := net.JoinHostPort("", strconv.Itoa(config.Config.Port))
	if err := app.Listen(addr, fiber.ListenConfig{
		EnablePrefork:     true,
		EnablePrintRoutes: true,
	}); err != nil {
		logger.Panic("Server failed to start", zap.Error(err))
	}

	if !fiber.IsChild() {
		logger.Info("Server exited gracefully")
	}
}
