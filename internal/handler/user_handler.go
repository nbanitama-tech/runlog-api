package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/internal/requestcontext"
	"github.com/nbanitama-tech/runlog-api/internal/usecase"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
	"github.com/nbanitama-tech/runlog-api/pkg/response"
	"github.com/nbanitama-tech/runlog-api/pkg/transport/dto"
)

// UserUseCase defines the interface for user-related use case operations, including registration and login. It abstracts the underlying implementation of user management, allowing the handler to interact with user data without being tightly coupled to a specific implementation.
type UserUseCase interface {
	Register(ctx context.Context, name, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (*usecase.LoginResult, error)
}

// UserHandler is responsible for handling user-related HTTP requests in the RunLog API application. It provides methods for user registration, login, and profile retrieval. The handler interacts with the UserUseCase interface to perform business logic operations related to user management.
type UserHandler struct {
	userUseCase UserUseCase
}

// NewUserHandler creates a new instance of UserHandler with the provided UserUseCase. It initializes the handler with the necessary dependencies to handle user-related HTTP requests, such as registration and login.
func NewUserHandler(userUseCase UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Create a new RunLog user account
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"Register request"
//	@Success		201		{object}	dto.UserResponseEnvelope
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	user, err := h.userUseCase.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, user)
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and return JWT
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login request"
//	@Success		200		{object}	dto.LoginResponseEnvelope
//	@Failure		401		{object}	dto.ErrorResponse
//	@Router			/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.userUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == pkgerrors.ErrInvalidCredentials {
			response.Unauthorized(c, "invalid email or password")
			return
		}
		response.InternalServerError(c)
		return

	}

	response.OK(c, gin.H{
		"token": result.Token,
		"user": gin.H{
			"id":    result.User.ID,
			"name":  result.User.Name,
			"email": result.User.Email,
		},
	})
}

// Profile godoc
//
//	@Summary		Get current user
//	@Description	Get authenticated user profile
//	@Tags			Users
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	dto.UserResponseEnvelope
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/users/profile [get]
func (h *UserHandler) Profile(c *gin.Context) {
	userID := requestcontext.UserID(c)
	email := requestcontext.Email(c)

	response.OK(c, gin.H{
		"id":    userID,
		"email": email,
	})
}
