package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nbanitama-tech/runlog-api/internal/requestcontext"
)

// RequestIDMiddleware is a Gin middleware function that generates or retrieves a unique request ID for each incoming HTTP request. It checks for the presence of an "X-Request-ID" header in the request; if absent, it generates a new UUID as the request ID. The middleware sets the request ID in the Gin context and includes it in the response headers, allowing for consistent tracking of requests throughout the application.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(requestcontext.RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
