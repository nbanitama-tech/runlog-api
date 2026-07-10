package usecase

import (
	"context"

	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/internal/repository"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
	"github.com/nbanitama-tech/runlog-api/pkg/infrastructure/auth"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase handles user-related business logic, including registration and login. It interacts with the UserRepository to perform operations on user data and utilizes JWT for authentication.
type UserUseCase struct {
	userRepo       repository.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

// LoginResult represents the result of a successful login operation, containing the generated JWT token and the authenticated user information.
type LoginResult struct {
	Token string
	User  *model.User
}

// NewUserUseCase creates a new instance of UserUseCase with the provided UserRepository, JWT secret, and JWT expiry hours. It initializes the use case with the necessary dependencies to handle user-related business logic, such as registration and login.
func NewUserUseCase(userRepo repository.UserRepository, jwtSecret string, jwtExpiryHours int) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

// Register registers a new user with the provided name, email, and password. It hashes the password using bcrypt and creates a new user in the UserRepository. If successful, it returns the created User object; otherwise, it returns an error.
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

// Login authenticates a user with the provided email and password. It retrieves the user from the UserRepository, compares the provided password with the stored password hash using bcrypt, and generates a JWT token if authentication is successful. If successful, it returns a LoginResult containing the token and user information; otherwise, it returns an error.
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
