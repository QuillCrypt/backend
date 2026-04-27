package config

import (
	"quillcrypt-backend/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

type config struct {
	Mode                int `default:"0"`
	Port                int `default:"8080"`
	Gh_ClientId         string
	Gh_ClientSecret     string
	Gh_Callback         string `default:"http://localhost:8080/auth/github/callback"`
	Google_ClientId     string
	Google_ClientSecret string
	Google_Callback     string `default:"http://localhost:8080/auth/google/callback"`
	RedisURL            string `default:""`
	PGURL               string `default:""`
	SessionSecret       string `default:"quillcrypt-secret-key"`
	LogFilePath         string `default:""`
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
		github.New(Config.Gh_ClientId, Config.Gh_ClientSecret, Config.Gh_Callback),
		google.New(Config.Google_ClientId, Config.Google_ClientSecret, Config.Google_Callback),
	)
}
