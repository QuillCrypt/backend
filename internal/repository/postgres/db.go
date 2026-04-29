package postgres

import (
	"context"
	"errors"
	"quillcrypt-backend/internal/config"
	"quillcrypt-backend/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *pgxpool.Pool

func InitDB() {
	if config.Config.PGURL == "" {
		logger.Panic("QC_PGURL missing in .env but is required.")
	}

	config, err := pgxpool.ParseConfig(config.Config.PGURL)
	if err != nil {
		logger.Panic("Unable to parse Postgres URL", zap.Error(err))
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Panic("Unable to create Postgres connection pool", zap.Error(err))
	}

	if err := DB.Ping(context.Background()); err != nil {
		logger.Panic("Unable to ping Postgres", zap.Error(err))
	}

	logger.Info("Postgres connection pool initialized")

	if !fiber.IsChild() {
		Migrate()
	}
}

func Migrate() {
	m, err := migrate.New(
		"file://migrations",
		config.Config.PGURL,
	)
	if err != nil {
		logger.Panic("Unable to create migrate instance", zap.Error(err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Panic("Unable to run migrations", zap.Error(err))
	}

	logger.Info("Postgres migrations applied successfully")
}
