package dto

// ErrorBody represents the structure of an error response body. It includes a code and a message that provide information about the error encountered during an API request. The struct is used to standardize error responses in the RunLog API application.
type ErrorBody struct {
	Code    string `json:"code" example:"UNAUTHORIZED"`
	Message string `json:"message" example:"invalid email or password"`
}

// ErrorResponse represents the structure of an error response. It includes a success flag set to false and an ErrorBody that contains details about the error. The struct is used to standardize error responses in the RunLog API application.
type ErrorResponse struct {
	Success bool      `json:"success" example:"false"`
	Error   ErrorBody `json:"error"`
}

// SuccessResponse represents the structure of a successful response. It includes a success flag set to true and an optional Data field that can hold any type of data. The struct is used to standardize successful responses in the RunLog API application.
type SuccessResponse struct {
	Success bool `json:"success" example:"true"`
}
