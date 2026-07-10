// Package requestcontext provides utility functions to extract specific values from the Gin context, such as request ID and user ID. These functions help in maintaining a consistent way to access these values across different parts of the application.
package requestcontext

import "github.com/gin-gonic/gin"

const (
	// RequestIDKey is the key used to store and retrieve the request ID from the Gin context. It is used to uniquely identify each incoming request for logging and tracing purposes.
	RequestIDKey = "request_id"
	// UserIDKey is the key used to store and retrieve the user ID from the Gin context. It is used to identify the authenticated user making the request, allowing for user-specific operations and access control.
	UserIDKey = "user_id"
)

// RequestID retrieves the request ID from the Gin context. It returns the request ID as a string, which can be used for logging, tracing, and correlating requests across different parts of the application.
func RequestID(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}

// UserID retrieves the user ID from the Gin context. It returns the user ID as a string, which can be used to identify the authenticated user making the request and perform user-specific operations and access control.
func UserID(c *gin.Context) string {
	return c.GetString(UserIDKey)
}
