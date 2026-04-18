package config

import (
	"quillcrypt-backend/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Mode int `default:"0"`
	Port int `default:"8080"`
}

var Config config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		logger.Panic("Cannot load .env")
	}
	err = envconfig.Process("QC", &Config)
	if err != nil {
		logger.Panic("Cannot process .env")
	}
}
