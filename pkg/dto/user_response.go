package dto

import "time"

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponseEnvelope struct {
	Success bool         `json:"success" example:"true"`
	Data    UserResponse `json:"data"`
}

type LoginResponseEnvelope struct {
	Success bool          `json:"success" example:"true"`
	Data    LoginResponse `json:"data"`
}
