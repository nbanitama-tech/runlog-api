package dto

type ErrorBody struct {
	Code    string `json:"code" example:"UNAUTHORIZED"`
	Message string `json:"message" example:"invalid email or password"`
}

type ErrorResponse struct {
	Success bool      `json:"success" example:"false"`
	Error   ErrorBody `json:"error"`
}

type SuccessResponse struct {
	Success bool `json:"success" example:"true"`
}
