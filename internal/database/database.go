package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(logger *slog.Logger) *pgxpool.Pool {
	dbDSN := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		os.Getenv("DB_CONNECTION"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), dbDSN)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		logger.Error("Failed to established connection", "error", err)
		os.Exit(1)
	}

	return pool
}
