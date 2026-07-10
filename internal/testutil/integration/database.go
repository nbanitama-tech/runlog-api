package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultTestDatabaseURL = "postgres://runlog:runlog_password@localhost:5432/runlog_test_db?sslmode=disable"

// NewDB creates a new PostgreSQL connection pool for integration tests. It checks the environment variable "RUN_INTEGRATION_TESTS" to determine whether to run the integration tests. If the variable is not set to "true", the test is skipped. The function retrieves the database URL from the "TEST_DATABASE_URL" environment variable or uses a default value if not provided. It returns a pointer to the pgxpool.Pool instance and any error encountered during the connection process.
func NewDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test")
	}

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = defaultTestDatabaseURL
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to connect test database: %v", err)
	}

	return db
}
