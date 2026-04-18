package main

import (
	"quillcrypt-backend/internal/config"
	"quillcrypt-backend/pkg/logger"

	"go.uber.org/zap"
)

func init() {
	config.LoadConfig()
}

func main() {
	logger.Init(config.Config.Mode)
	defer logger.Log.Sync()
	
	logger.Info("Starting QuillCrypt backend server...",
		zap.Int("Port", config.Config.Port),
	)
}