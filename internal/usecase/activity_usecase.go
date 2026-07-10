// Package usecase provides the business logic layer for the RunLog API application. It defines the ActivityUseCase struct, which implements the core functionality for managing user activities, including creating, listing, retrieving, updating, and deleting activities. The use case interacts with the ActivityRepository interface to perform data access operations and encapsulates the business rules and calculations related to user activities.
package usecase

import (
	"context"
	"time"

	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/internal/repository"
)

// ActivityUseCase is responsible for managing user activities in the RunLog API application. It provides methods for creating, listing, retrieving, updating, and deleting activities. The use case interacts with the ActivityRepository interface to perform data access operations and encapsulates the business rules and calculations related to user activities.
type ActivityUseCase struct {
	activityRepo repository.ActivityRepository
}

// NewActivityUseCase creates a new instance of ActivityUseCase with the provided ActivityRepository. It initializes the use case with the necessary dependencies to manage user activities, allowing it to perform business logic operations related to activity management.
func NewActivityUseCase(activityRepo repository.ActivityRepository) *ActivityUseCase {
	return &ActivityUseCase{activityRepo: activityRepo}
}

// Create creates a new activity for a user based on the provided parameters. It calculates the average pace based on the distance and duration, and sets default values for sport type if not provided. The method interacts with the ActivityRepository to persist the activity in the database and returns the created Activity object. If an error occurs during the creation process, it returns the corresponding error.
func (u *ActivityUseCase) Create(
	ctx context.Context,
	userID string,
	title string,
	sportType string,
	distanceKM float64,
	durationSeconds int,
	elevationGainM int,
	activityDate time.Time,
	notes string,
) (*model.Activity, error) {
	avgPaceSeconds := 0
	if distanceKM > 0 {
		avgPaceSeconds = int(float64(durationSeconds) / distanceKM)
	}

	if sportType == "" {
		sportType = "running"
	}

	activity := &model.Activity{
		UserID:          userID,
		Title:           title,
		SportType:       sportType,
		DistanceKM:      distanceKM,
		DurationSeconds: durationSeconds,
		AvgPaceSeconds:  avgPaceSeconds,
		ElevationGainM:  elevationGainM,
		ActivityDate:    activityDate,
		Notes:           notes,
	}

	if err := u.activityRepo.Create(ctx, activity); err != nil {
		return nil, err
	}

	return activity, nil
}

// ListByUserID retrieves a list of activities for a specific user based on the provided filter criteria. It interacts with the ActivityRepository to fetch the activities from the database and returns an ActivityListResult containing the activities and the total count. The filter parameter allows for pagination, sport type filtering, date range filtering, and sorting options.
func (u *ActivityUseCase) ListByUserID(ctx context.Context, userID string, filter model.ActivityFilter) (*model.ActivityListResult, error) {
	return u.activityRepo.FindByUserID(ctx, userID, filter)
}

// GetByID retrieves a specific activity for a user based on the provided user ID and activity ID. It interacts with the ActivityRepository to fetch the activity from the database and returns the corresponding Activity object. If the activity is not found, it returns an error indicating that the activity does not exist.
func (u *ActivityUseCase) GetByID(ctx context.Context, userID, activityID string) (*model.Activity, error) {
	return u.activityRepo.FindByID(ctx, userID, activityID)
}

// Update modifies an existing activity for a user based on the provided parameters. It calculates the average pace based on the distance and duration, and updates the activity's attributes accordingly. The method interacts with the ActivityRepository to persist the changes in the database and returns the updated Activity object. If an error occurs during the update process, it returns the corresponding error.
func (u *ActivityUseCase) Update(
	ctx context.Context,
	userID string,
	activityID string,
	title string,
	sportType string,
	distanceKM float64,
	durationSeconds int,
	elevationGainM int,
	activityDate time.Time,
	notes string,
) (*model.Activity, error) {
	avgPaceSeconds := 0
	if distanceKM > 0 {
		avgPaceSeconds = int(float64(durationSeconds) / distanceKM)
	}

	if sportType == "" {
		sportType = "running"
	}

	activity := &model.Activity{
		ID:              activityID,
		UserID:          userID,
		Title:           title,
		SportType:       sportType,
		DistanceKM:      distanceKM,
		DurationSeconds: durationSeconds,
		AvgPaceSeconds:  avgPaceSeconds,
		ElevationGainM:  elevationGainM,
		ActivityDate:    activityDate,
		Notes:           notes,
	}

	if err := u.activityRepo.Update(ctx, activity); err != nil {
		return nil, err
	}

	return activity, nil
}

// Delete removes a specific activity for a user based on the provided user ID and activity ID. It interacts with the ActivityRepository to delete the activity from the database. If an error occurs during the deletion process, it returns the corresponding error.
func (u *ActivityUseCase) Delete(ctx context.Context, userID, activityID string) error {
	return u.activityRepo.Delete(ctx, userID, activityID)
}
