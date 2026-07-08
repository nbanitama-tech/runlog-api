package dto

type CreateActivityRequest struct {
	Title           string  `json:"title" binding:"required"`
	SportType       string  `json:"sport_type"`
	DistanceKM      float64 `json:"distance_km" binding:"required,gt=0"`
	DurationSeconds int     `json:"duration_seconds" binding:"required,gt=0"`
	ElevationGainM  int     `json:"elevation_gain_m"`
	ActivityDate    string  `json:"activity_date" binding:"required"`
	Notes           string  `json:"notes"`
}

type UpdateActivityRequest struct {
	Title           string  `json:"title" binding:"required"`
	SportType       string  `json:"sport_type"`
	DistanceKM      float64 `json:"distance_km" binding:"required,gt=0"`
	DurationSeconds int     `json:"duration_seconds" binding:"required,gt=0"`
	ElevationGainM  int     `json:"elevation_gain_m"`
	ActivityDate    string  `json:"activity_date" binding:"required"`
	Notes           string  `json:"notes"`
}
