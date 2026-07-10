package model

import "time"

// User represents a user in the system with attributes such as ID, Name, Email, PasswordHash, CreatedAt, and UpdatedAt. It serves as the core data structure for managing user information in the RunLog API application.
type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
