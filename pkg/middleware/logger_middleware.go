package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware is a Gin middleware function that logs HTTP requests and responses. It captures details such as request method, path, status code, latency, client IP, user agent, and user ID (if available). The log messages are written using the provided slog.Logger instance, allowing for structured logging in JSON format. This middleware is useful for monitoring and debugging API requests in the RunLog API application.
func LoggerMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		log.Info("http_request",
			"request_id", c.GetString("request_id"),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", latency.Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"user_id", c.GetString("user_id"),
		)
	}
}
