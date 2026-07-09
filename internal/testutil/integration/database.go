package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultTestDatabaseURL = "postgres://runlog:runlog_password@localhost:5432/runlog_test_db?sslmode=disable"

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
