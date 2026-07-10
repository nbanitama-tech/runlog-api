package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/requestcontext"
)

// LoggerMiddleware is a Gin middleware function that logs HTTP requests and responses. It captures details such as request method, path, status code, latency, client IP, user agent, and user ID (if available). The log messages are written using the provided slog.Logger instance, allowing for structured logging in JSON format. This middleware is useful for monitoring and debugging API requests in the RunLog API application.
func LoggerMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		log.Info(
			"http_request",

			slog.String("request_id", requestcontext.RequestID(c)),
			slog.String("user_id", requestcontext.UserID(c)),

			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),

			slog.Int("status", c.Writer.Status()),

			slog.Duration("latency", latency),

			slog.String("client_ip", c.ClientIP()),

			slog.String("user_agent", c.Request.UserAgent()),
		)
	}
}
