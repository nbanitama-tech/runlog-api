// Package integration provides utility functions for integration testing in the RunLog API application. It includes the Cleanup function, which is used to reset the state of the test database by truncating relevant tables and restarting their identity sequences. This ensures that each test runs in a clean environment without interference from previous tests.
package integration

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Cleanup resets the state of the test database by truncating the "activities" and "users" tables and restarting their identity sequences. It is intended to be used in integration tests to ensure that each test runs in a clean environment without interference from previous tests. The function takes a testing.T instance and a PostgreSQL connection pool as parameters, allowing it to execute the necessary SQL commands to perform the cleanup operation.
func Cleanup(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	_, err := db.Exec(
		context.Background(),
		`
		TRUNCATE TABLE
			activities,
			users
		RESTART IDENTITY CASCADE;
		`,
	)
	if err != nil {
		t.Fatalf("failed to cleanup test database: %v", err)
	}
}
