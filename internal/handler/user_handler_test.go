package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/internal/usecase"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
)

type mockUserUseCase struct {
	registerUser *model.User
	loginResult  *usecase.LoginResult
	err          error
}

func (m *mockUserUseCase) Register(_ context.Context, _, _, _ string) (*model.User, error) {
	return m.registerUser, m.err
}

func (m *mockUserUseCase) Login(_ context.Context, _, _ string) (*usecase.LoginResult, error) {
	return m.loginResult, m.err
}

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockUserUseCase{
		registerUser: &model.User{
			ID:        "user-123",
			Name:      "Novandi",
			Email:     "novandi@example.com",
			CreatedAt: time.Now(),
		},
	}

	h := NewUserHandler(mockUC)

	r := gin.New()
	r.POST("/users/register", h.Register)

	body := map[string]string{
		"name":     "Novandi",
		"email":    "novandi@example.com",
		"password": "password123",
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockUserUseCase{
		loginResult: &usecase.LoginResult{
			Token: "test-token",
			User: &model.User{
				ID:    "user-123",
				Name:  "Novandi",
				Email: "novandi@example.com",
			},
		},
	}

	h := NewUserHandler(mockUC)

	r := gin.New()
	r.POST("/users/login", h.Login)

	body := map[string]string{
		"email":    "novandi@example.com",
		"password": "password123",
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_Register_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockUserUseCase{}

	h := NewUserHandler(mockUC)

	r := gin.New()
	r.POST("/users/register", h.Register)

	body := map[string]string{
		"name":     "",
		"email":    "invalid-email",
		"password": "short",
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockUserUseCase{
		err: pkgerrors.ErrInvalidCredentials,
	}

	h := NewUserHandler(mockUC)

	r := gin.New()
	r.POST("/users/login", h.Login)

	body := map[string]string{
		"email":    "novandi@example.com",
		"password": "wrong-password",
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d, body: %s", w.Code, w.Body.String())
	}
}
