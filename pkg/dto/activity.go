// Package dto defines the data transfer objects (DTOs) used in the RunLog API application. It includes request and response structures for creating and updating user activities, as well as other related operations. The DTOs are used to facilitate communication between the API handlers and the underlying business logic, ensuring that data is properly structured and validated before being processed.
package dto

// CreateActivityRequest represents the request payload for creating a new activity. It includes fields such as Title, SportType, DistanceKM, DurationSeconds, ElevationGainM, ActivityDate, and Notes. The struct is used to validate and transfer data from the API request to the business logic layer.
type CreateActivityRequest struct {
	Title           string  `json:"title" binding:"required"`
	SportType       string  `json:"sport_type"`
	DistanceKM      float64 `json:"distance_km" binding:"required,gt=0"`
	DurationSeconds int     `json:"duration_seconds" binding:"required,gt=0"`
	ElevationGainM  int     `json:"elevation_gain_m"`
	ActivityDate    string  `json:"activity_date" binding:"required"`
	Notes           string  `json:"notes"`
}

// UpdateActivityRequest represents the request payload for updating an existing activity. It includes fields such as Title, SportType, DistanceKM, DurationSeconds, ElevationGainM, ActivityDate, and Notes. The struct is used to validate and transfer data from the API request to the business logic layer when updating an activity.
type UpdateActivityRequest struct {
	Title           string  `json:"title" binding:"required"`
	SportType       string  `json:"sport_type"`
	DistanceKM      float64 `json:"distance_km" binding:"required,gt=0"`
	DurationSeconds int     `json:"duration_seconds" binding:"required,gt=0"`
	ElevationGainM  int     `json:"elevation_gain_m"`
	ActivityDate    string  `json:"activity_date" binding:"required"`
	Notes           string  `json:"notes"`
}
