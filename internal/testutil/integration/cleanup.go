package integration

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
