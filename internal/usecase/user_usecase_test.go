package usecase

import (
	"context"
	"testing"

	"github.com/nbanitama-tech/runlog-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	createdUser *model.User
	findUser    *model.User
	err         error
}

func (m *mockUserRepository) Create(_ context.Context, user *model.User) error {
	m.createdUser = user
	return m.err
}

func (m *mockUserRepository) FindByEmail(_ context.Context, _ string) (*model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.findUser, nil
}

func TestUserUseCase_Register(t *testing.T) {
	repo := &mockUserRepository{}

	uc := NewUserUseCase(repo, "test-secret", 24)

	user, err := uc.Register(
		context.Background(),
		"Novandi",
		"novandi@example.com",
		"password123",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != "novandi@example.com" {
		t.Fatalf("expected email novandi@example.com, got %s", user.Email)
	}

	if user.PasswordHash == "password123" {
		t.Fatal("expected password to be hashed")
	}

	if repo.createdUser == nil {
		t.Fatal("expected user to be created")
	}
}

func TestUserUseCase_Login_Success(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repo := &mockUserRepository{
		findUser: &model.User{
			ID:           "user-123",
			Name:         "Novandi",
			Email:        "novandi@example.com",
			PasswordHash: string(passwordHash),
		},
	}

	uc := NewUserUseCase(repo, "test-secret", 24)

	result, err := uc.Login(
		context.Background(),
		"novandi@example.com",
		"password123",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Token == "" {
		t.Fatal("expected token to be generated")
	}

	if result.User.Email != "novandi@example.com" {
		t.Fatalf("expected email novandi@example.com, got %s", result.User.Email)
	}
}

func TestUserUseCase_Login_InvalidPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repo := &mockUserRepository{
		findUser: &model.User{
			ID:           "user-123",
			Name:         "Novandi",
			Email:        "novandi@example.com",
			PasswordHash: string(passwordHash),
		},
	}

	uc := NewUserUseCase(repo, "test-secret", 24)

	_, err = uc.Login(
		context.Background(),
		"novandi@example.com",
		"wrong-password",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
