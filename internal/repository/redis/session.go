package redis

import (
	"context"
	"quillcrypt-backend/internal/config"
	"quillcrypt-backend/pkg/logger"
	"runtime"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Store *session.Store
var Client *redis.Client
var Redsync *redsync.Redsync

// fiberStorage implements fiber.Storage interface using our shared Client
type fiberStorage struct{}

func (s fiberStorage) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	val, err := Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (s fiberStorage) Get(key string) ([]byte, error) {
	return s.GetWithContext(context.Background(), key)
}

func (s fiberStorage) SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error {
	return Client.Set(ctx, key, val, exp).Err()
}

func (s fiberStorage) Set(key string, val []byte, exp time.Duration) error {
	return s.SetWithContext(context.Background(), key, val, exp)
}

func (s fiberStorage) DeleteWithContext(ctx context.Context, key string) error {
	return Client.Del(ctx, key).Err()
}

func (s fiberStorage) Delete(key string) error {
	return s.DeleteWithContext(context.Background(), key)
}

func (s fiberStorage) ResetWithContext(ctx context.Context) error {
	return Client.FlushDB(ctx).Err()
}

func (s fiberStorage) Reset() error {
	return s.ResetWithContext(context.Background())
}

func (s fiberStorage) Close() error {
	return nil // Managed by Client
}

func InitSession() {
	if config.Config.RedisURL == "" {
		logger.Panic("QC_REDISURL missing in .env but is required.")
	}

	opts, err := redis.ParseURL(config.Config.RedisURL)
	if err != nil {
		logger.Panic("Unable to parse Redis URL", zap.Error(err))
	}
	opts.PoolSize = 10 * runtime.GOMAXPROCS(0)
	Client = redis.NewClient(opts)

	Store = session.NewStore(session.Config{
		Storage:   fiberStorage{},
		Extractor: extractors.FromCookie("qc_session"),
	})

	// Setup Redsync using the same Client
	pool := goredis.NewPool(Client)
	Redsync = redsync.New(pool)

	logger.Info("Redis session storage initialized")
}
