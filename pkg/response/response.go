// Package response provides utility functions for sending standardized HTTP responses in the RunLog API application. It includes functions for sending success and error responses with appropriate HTTP status codes and JSON payloads. The package defines response structures for success and error cases, allowing consistent response formatting across the application.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorBody represents the structure of an error response body. It includes a code and a message that provide information about the error encountered during an API request. The struct is used to standardize error responses in the RunLog API application.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse represents the structure of a successful response. It includes a success flag set to true and an optional Data field that can hold any type of data. The struct is used to standardize successful responses in the RunLog API application.
type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data,omitempty"`
}

// ErrorResponse struct represents the structure of an error response. It includes a success flag set to false and an ErrorBody that contains details about the error. The struct is used to standardize error responses in the RunLog API application.
type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

// OK sends a successful response with HTTP status 200 (OK) and the provided data.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// Created sends a successful response with HTTP status 201 (Created) and the provided data.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a successful response with HTTP status 204 (No Content).
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends an error response with HTTP status 400 (Bad Request) and the provided error message.
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error: ErrorBody{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	})
}

// Unauthorized sends an error response with HTTP status 401 (Unauthorized) and the provided error message.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Error: ErrorBody{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
	})
}

// NotFound sends an error response with HTTP status 404 (Not Found) and the provided error message.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Success: false,
		Error: ErrorBody{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

// InternalServerError sends an error response with HTTP status 500 (Internal Server Error) and a generic error message.
func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error: ErrorBody{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		},
	})
}
