package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/usecase"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.userUseCase.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := h.userUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(200, gin.H{
		"token": result.Token,
		"user": gin.H{
			"id":    result.User.ID,
			"name":  result.User.Name,
			"email": result.User.Email,
		},
	})
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID := c.GetString("user_id")
	email := c.GetString("email")

	c.JSON(http.StatusOK, gin.H{
		"id":    userID,
		"email": email,
	})
}
