package dto

// RegisterRequest represents the request payload for user registration. It includes fields for Name, Email, and Password, with validation rules specified using struct tags. The struct is used to capture and validate user input when creating a new user account in the RunLog API application.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the request payload for user login. It includes fields for Email and Password, with validation rules specified using struct tags. The struct is used to capture and validate user input when authenticating a user in the RunLog API application.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
