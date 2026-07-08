package usecase

import (
	"context"

	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/internal/repository"
	"github.com/nbanitama-tech/runlog-api/pkg/auth"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo       repository.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

type LoginResult struct {
	Token string
	User  *model.User
}

func NewUserUseCase(userRepo repository.UserRepository, jwtSecret string, jwtExpiryHours int) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

func (u *UserUseCase) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, pkgerrors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, pkgerrors.ErrInvalidCredentials
	}

	expiryHours := u.jwtExpiryHours
	if expiryHours <= 0 {
		expiryHours = 24 // Default to 24 hours if not set
	}

	token, err := auth.GenerateToken(user.ID, user.Email, u.jwtSecret, expiryHours)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
		User:  user,
	}, nil
}
