package usecase

import (
	"context"
	"time"

	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/internal/repository"
)

type ActivityUseCase struct {
	activityRepo *repository.ActivityRepository
}

func NewActivityUseCase(activityRepo *repository.ActivityRepository) *ActivityUseCase {
	return &ActivityUseCase{activityRepo: activityRepo}
}

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

func (u *ActivityUseCase) ListByUserID(ctx context.Context, userID string) ([]model.Activity, error) {
	return u.activityRepo.FindByUserID(ctx, userID)
}

func (u *ActivityUseCase) GetByID(ctx context.Context, userID, activityID string) (*model.Activity, error) {
	return u.activityRepo.FindByID(ctx, userID, activityID)
}

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
