// Package model defines the data structures used in the RunLog API application. It includes the Activity struct, which represents a user's activity with fields such as ID, UserID, Title, SportType, DistanceKM, DurationSeconds, AvgPaceSeconds, ElevationGainM, ActivityDate, Notes, CreatedAt, and UpdatedAt. These models are used throughout the application for data representation and manipulation.
package model

import "time"

// Activity represents a user's activity with various attributes such as ID, UserID, Title, SportType, DistanceKM, DurationSeconds, AvgPaceSeconds, ElevationGainM, ActivityDate, Notes, CreatedAt, and UpdatedAt. It serves as the core data structure for managing user activities in the RunLog API application.
type Activity struct {
	ID              string
	UserID          string
	Title           string
	SportType       string
	DistanceKM      float64
	DurationSeconds int
	AvgPaceSeconds  int
	ElevationGainM  int
	ActivityDate    time.Time
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
