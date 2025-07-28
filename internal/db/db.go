package db

import (
	"auth-service/internal/config"
	"auth-service/internal/logger"

	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Log.Errorf("failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return pool, nil
}
