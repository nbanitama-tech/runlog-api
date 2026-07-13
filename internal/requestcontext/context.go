// Package requestcontext provides utility functions to extract specific values from the Gin context, such as request ID, user ID, and email. These functions help in retrieving context-specific information that may be set during the request lifecycle, facilitating better request tracking and user identification within the application.
package requestcontext

import "github.com/gin-gonic/gin"

const (
	// RequestIDKey is the key used to store and retrieve the request ID from the Gin context. It is typically set by middleware that generates a unique request ID for each incoming HTTP request, allowing for better tracking and correlation of requests in logs and error reporting.
	RequestIDKey = "request_id"
	// UserIDKey is the key used to store and retrieve the user ID from the Gin context. It is typically set after successful authentication, allowing handlers and middleware to access the authenticated user's ID for authorization and personalization purposes.
	UserIDKey = "user_id"
	// EmailKey is the key used to store and retrieve the user's email address from the Gin context. It is typically set after successful authentication, allowing handlers and middleware to access the authenticated user's email for identification and communication purposes.
	EmailKey = "email"
)

// RequestID retrieves the request ID from the Gin context. It returns the request ID as a string, which can be used for logging, tracing, and correlating requests throughout the application. If the request ID is not set in the context, it returns an empty string.
func RequestID(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}

// UserID retrieves the user ID from the Gin context. It returns the user ID as a string, which can be used for authorization, personalization, and tracking user-specific actions within the application. If the user ID is not set in the context, it returns an empty string.
func UserID(c *gin.Context) string {
	return c.GetString(UserIDKey)
}

// Email retrieves the user's email address from the Gin context. It returns the email as a string, which can be used for identification, communication, and personalization purposes within the application. If the email is not set in the context, it returns an empty string.
func Email(c *gin.Context) string {
	return c.GetString(EmailKey)
}
