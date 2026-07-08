package model

import "time"

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
