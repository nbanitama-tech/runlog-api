package dto

import "time"

// ActivityResponse represents the response payload for an activity. It includes fields such as ID, Title, SportType, DistanceKM, DurationSeconds, AvgPaceSeconds, ElevationGainM, ActivityDate, Notes, CreatedAt, and UpdatedAt. The struct is used to transfer activity data from the business logic layer to the API response.
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

// ActivityResponseEnvelope represents the response envelope for a single activity. It includes a success flag and the activity data. The struct is used to wrap the activity response in a consistent format for API responses.
type ActivityResponseEnvelope struct {
	Success bool             `json:"success" example:"true"`
	Data    ActivityResponse `json:"data"`
}

// ActivityListResponseEnvelope represents the response envelope for a list of activities. It includes a success flag and a slice of activity data. The struct is used to wrap the list of activity responses in a consistent format for API responses.
type ActivityListResponseEnvelope struct {
	Success bool               `json:"success" example:"true"`
	Data    []ActivityResponse `json:"data"`
}
