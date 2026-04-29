package config

import (
	"quillcrypt-backend/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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
	MobileCallback  string `default:"quillcrypt://callback"`
}

var Config config
var OAuth2Config *oauth2.Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found")
	}
	err = envconfig.Process("QC", &Config)
	if err != nil {
		logger.Panic("Cannot process env vars")
	}

	OAuth2Config = &oauth2.Config{
		ClientID:     Config.Gh_ClientId,
		ClientSecret: Config.Gh_ClientSecret,
		RedirectURL:  Config.Gh_Callback,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"read:user", "user:email"},
	}
}
