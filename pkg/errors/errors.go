// Package errors defines custom error variables for the RunLog API application. These errors represent specific failure scenarios that can occur during user and activity operations, such as user not found, invalid credentials, email already exists, and activity not found. The package uses the standard library's errors package to create these error variables, which can be used throughout the application for consistent error handling and messaging.
package errors

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found in the database.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials is returned when the provided email or password is invalid during user authentication.
	ErrInvalidCredentials = errors.New("invalid email or password")
	// ErrEmailAlreadyExists is returned when attempting to register a user with an email that already exists in the database.
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrActivityNotFound is returned when an activity is not found in the database.
	ErrActivityNotFound = errors.New("activity not found")
)
