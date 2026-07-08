package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfig struct {
	URL string
}

func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, databaseURL)
}
