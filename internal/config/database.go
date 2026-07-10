package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseConfig holds the configuration settings for the database connection in the RunLog API application. It defines the URL used to connect to the database, allowing the application to establish a connection with the specified database server.
type DatabaseConfig struct {
	URL string
}

// NewPostgresPool creates a new PostgreSQL connection pool using the provided database URL. It returns a pointer to the pgxpool.Pool instance and any error encountered during the connection process. The function uses the pgxpool package to manage the connection pool, allowing efficient handling of multiple database connections in the RunLog API application.
func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, databaseURL)
}
