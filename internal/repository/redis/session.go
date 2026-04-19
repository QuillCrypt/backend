package redis

import (
	"net/url"
	"quillcrypt-backend/internal/config"
	"quillcrypt-backend/pkg/logger"
	"runtime"

	"github.com/boj/redistore"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

	var Store *session.Store

	func InitSession() {
	var storage *redis.Storage
	if config.Config.RedisURL == "" {
		logger.Panic("QC_REDISURL missing in .env but is required.")
	}
	storage = redis.New(redis.Config{
		URL: config.Config.RedisURL,
		Database: 0,
		Reset: false,
		PoolSize: 10 * runtime.GOMAXPROCS(0),
	})

	Store = session.NewStore(session.Config{
		Storage: storage,
		Extractor: extractors.FromCookie("qc_session"),
	})

	// Setup Gothic Redis Store
	u, err := url.Parse(config.Config.RedisURL)
	if err != nil {
		logger.Panic("Unable to parse Redis URL for Gothic", zap.Error(err))
	}
	host := u.Host
	if host == "" {
		host = "localhost:6379"
	}

	gothStore, err := redistore.NewRediStore(10, "tcp", host, "", "", []byte(config.Config.SessionSecret))
	if err != nil {
		logger.Panic("Unable to initialize Gothic Redis store", zap.Error(err))
	}
	gothic.Store = gothStore

	logger.Info("Redis session storage initialized")
	}


