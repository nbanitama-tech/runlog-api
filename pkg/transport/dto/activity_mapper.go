package dto

import "github.com/nbanitama-tech/runlog-api/internal/model"

// ToActivityResponse converts a model.Activity to an ActivityResponse DTO. It maps the fields from the model to the response structure, formatting the ActivityDate as a string in "YYYY-MM-DD" format. This function is used to prepare activity data for API responses, ensuring that the data is properly structured and formatted for client consumption.
func ToActivityResponse(activity model.Activity) ActivityResponse {
	return ActivityResponse{
		ID:              activity.ID,
		Title:           activity.Title,
		SportType:       activity.SportType,
		DistanceKM:      activity.DistanceKM,
		DurationSeconds: activity.DurationSeconds,
		AvgPaceSeconds:  activity.AvgPaceSeconds,
		ElevationGainM:  activity.ElevationGainM,
		ActivityDate:    activity.ActivityDate.Format("2006-01-02"),
		Notes:           activity.Notes,
		CreatedAt:       activity.CreatedAt,
		UpdatedAt:       activity.UpdatedAt,
	}
}

// ToActivityResponses converts a slice of model.Activity to a slice of ActivityResponse DTOs. It iterates over the input activities, converting each one using the ToActivityResponse function, and returns a new slice containing the corresponding ActivityResponse objects. This function is useful for preparing lists of activity data for API responses, ensuring that the data is properly structured and formatted for client consumption.
func ToActivityResponses(activities []model.Activity) []ActivityResponse {
	responses := make([]ActivityResponse, 0, len(activities))

	for _, activity := range activities {
		responses = append(responses, ToActivityResponse(activity))
	}

	return responses
}
