package config

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseConfig holds the configuration settings for the database connection in the RunLog API application. It defines the URL used to connect to the database, allowing the application to establish a connection with the specified database server.
type DatabaseConfig struct {
	URL string
}

// NewPostgresPool creates a new PostgreSQL connection pool using the provided database URL and query tracer. It parses the database URL, configures the connection settings, and returns a pgxpool.Pool instance that can be used to manage database connections efficiently. The function also handles any errors that may occur during the configuration process.
func NewPostgresPool(
	ctx context.Context,
	databaseURL string,
	tracer pgx.QueryTracer,
) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.Tracer = tracer

	return pgxpool.NewWithConfig(ctx, poolConfig)
}
