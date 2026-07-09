package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/usecase"
	"github.com/nbanitama-tech/runlog-api/pkg/dto"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
	"github.com/nbanitama-tech/runlog-api/pkg/response"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
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
	userID := c.GetString("user_id")
	email := c.GetString("email")

	response.OK(c, gin.H{
		"id":    userID,
		"email": email,
	})
}
