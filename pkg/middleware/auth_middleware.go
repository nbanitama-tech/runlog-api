// Package middleware provides middleware functions for the RunLog API application. It includes the AuthMiddleware function, which is responsible for validating JWT tokens in incoming HTTP requests. The middleware checks for the presence of an Authorization header, verifies the token's validity, and extracts user information from the token claims. If the token is invalid or missing, the middleware responds with an appropriate HTTP error status and message, preventing unauthorized access to protected routes.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/pkg/auth"
)

// AuthMiddleware is a Gin middleware function that validates JWT tokens in incoming HTTP requests. It checks for the presence of an Authorization header, verifies the token's validity using the provided JWT secret, and extracts user information from the token claims. If the token is invalid or missing, it responds with an HTTP 401 Unauthorized status and an error message, preventing unauthorized access to protected routes.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header must use Bearer token",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := auth.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
