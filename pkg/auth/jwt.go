// Package auth provides JWT authentication functionality for the RunLog API application. It includes functions to generate and validate JWT tokens, as well as a Claims struct that represents the user information stored in the token. The package uses the github.com/golang-jwt/jwt/v5 library for JWT handling and supports token expiration and validation.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the user information stored in the JWT token. It includes the UserID and Email fields, which are used to identify the authenticated user in the RunLog API application.
type Claims struct {
	UserID string
	Email  string
}

// GenerateToken generates a JWT token for the given user ID and email using the provided secret and expiry duration in hours. It creates a new token with the specified claims and signs it using the HS256 signing method. The generated token can be used for authenticating requests in the RunLog API application.
func GenerateToken(userID, email, secret string, expiryHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken validates a JWT token against the provided secret and returns the parsed claims if valid. It uses the HS256 signing method to verify the token's integrity. If the token is valid, it returns the Claims struct containing the user information; otherwise, it returns an error.
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id claim")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email claim")
	}

	return &Claims{
		UserID: userID,
		Email:  email,
	}, nil
}
