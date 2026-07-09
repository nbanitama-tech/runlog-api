package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/internal/testutil/integration"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
)

func TestUserRepository_CreateAndFindByEmail(t *testing.T) {
	db := integration.NewDB(t)
	defer db.Close()

	integration.Cleanup(t, db)
	defer integration.Cleanup(t, db)

	ctx := context.Background()
	repo := NewUserRepository(db)

	user := &model.User{
		Name:         "Novandi",
		Email:        "novandi.integration@example.com",
		PasswordHash: "hashed-password",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == "" {
		t.Fatal("expected user ID to be generated")
	}

	foundUser, err := repo.FindByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if foundUser.Email != user.Email {
		t.Fatalf("expected email %s, got %s", user.Email, foundUser.Email)
	}

	if foundUser.Name != user.Name {
		t.Fatalf("expected name %s, got %s", user.Name, foundUser.Name)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := integration.NewDB(t)
	defer db.Close()

	integration.Cleanup(t, db)
	defer integration.Cleanup(t, db)

	ctx := context.Background()
	repo := NewUserRepository(db)

	_, err := repo.FindByEmail(ctx, "notfound@example.com")

	if !errors.Is(err, pkgerrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
