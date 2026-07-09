package dto

import "time"

type ActivityResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	SportType       string    `json:"sport_type"`
	DistanceKM      float64   `json:"distance_km"`
	DurationSeconds int       `json:"duration_seconds"`
	AvgPaceSeconds  int       `json:"avg_pace_seconds"`
	ElevationGainM  int       `json:"elevation_gain_m"`
	ActivityDate    string    `json:"activity_date"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

type ActivityResponseEnvelope struct {
	Success bool             `json:"success" example:"true"`
	Data    ActivityResponse `json:"data"`
}

type ActivityListResponseEnvelope struct {
	Success bool               `json:"success" example:"true"`
	Data    []ActivityResponse `json:"data"`
}
