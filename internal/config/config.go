package config

import (
	"quillcrypt-backend/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
)

type config struct {
	Mode            int `default:"0"`
	Port            int `default:"8080"`
	Gh_ClientId     string
	Gh_ClientSecret string
	Gh_Callback     string `default:"http://localhost:8080/auth/github/callback"`
	RedisURL        string
	PGURL           string
	SessionSecret   string `default:"quillcrypt-secret-key"`
	LogFilePath     string
}

var Config config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found")
	}
	err = envconfig.Process("QC", &Config)
	if err != nil {
		logger.Panic("Cannot process env vars")
	}

	goth.UseProviders(
		github.New(Config.Gh_ClientId, Config.Gh_ClientSecret, Config.Gh_Callback, "read:user", "user:email"),
	)
}
