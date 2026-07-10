package dto

import "time"

// UserResponse represents the response payload for user information. It includes fields such as ID, Name, Email, and CreatedAt. The struct is used to transfer user data from the business logic layer to the API response in the RunLog API application.
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// LoginResponse represents the response payload for user login. It includes a JWT token and user details. The struct is used to transfer authentication information from the business logic layer to the API response in the RunLog API application.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponseEnvelope represents the response envelope for user-related operations. It includes a success flag and a UserResponse that contains user details. The struct is used to wrap the user response in a consistent format for API responses in the RunLog API application.
type UserResponseEnvelope struct {
	Success bool         `json:"success" example:"true"`
	Data    UserResponse `json:"data"`
}

// LoginResponseEnvelope represents the response envelope for user login operations. It includes a success flag and a LoginResponse that contains authentication information. The struct is used to wrap the login response in a consistent format for API responses in the RunLog API application.
type LoginResponseEnvelope struct {
	Success bool          `json:"success" example:"true"`
	Data    LoginResponse `json:"data"`
}
