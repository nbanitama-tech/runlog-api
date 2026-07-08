package dto

import "github.com/nbanitama-tech/runlog-api/internal/model"

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

func ToActivityResponses(activities []model.Activity) []ActivityResponse {
	responses := make([]ActivityResponse, 0, len(activities))

	for _, activity := range activities {
		responses = append(responses, ToActivityResponse(activity))
	}

	return responses
}
