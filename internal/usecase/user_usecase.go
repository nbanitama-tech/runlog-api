package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/internal/repository"
	"github.com/nbanitama-tech/runlog-api/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo       *repository.UserRepository
	jwtSecret      string
	jwtExpiryHours string
}

type LoginResult struct {
	Token string
	User  *model.User
}

func NewUserUseCase(userRepo *repository.UserRepository, jwtSecret, jwtExpiryHours string) *UserUseCase {
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
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	expiryHours, err := strconv.Atoi(u.jwtExpiryHours)
	if err != nil {
		expiryHours = 24
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
