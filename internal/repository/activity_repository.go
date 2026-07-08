package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nbanitama-tech/runlog-api/internal/model"
)

type ActivityRepository struct {
	db *pgxpool.Pool
}

func NewActivityRepository(db *pgxpool.Pool) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(ctx context.Context, activity *model.Activity) error {
	query := `
		INSERT INTO activities (
			user_id, title, sport_type, distance_km, duration_seconds,
			avg_pace_seconds, elevation_gain_m, activity_date, notes
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		activity.UserID,
		activity.Title,
		activity.SportType,
		activity.DistanceKM,
		activity.DurationSeconds,
		activity.AvgPaceSeconds,
		activity.ElevationGainM,
		activity.ActivityDate,
		activity.Notes,
	).Scan(&activity.ID, &activity.CreatedAt, &activity.UpdatedAt)
}

func (r *ActivityRepository) FindByUserID(ctx context.Context, userID string) ([]model.Activity, error) {
	query := `
		SELECT id, user_id, title, sport_type, distance_km, duration_seconds,
		       avg_pace_seconds, elevation_gain_m, activity_date, notes,
		       created_at, updated_at
		FROM activities
		WHERE user_id = $1
		ORDER BY activity_date DESC, created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []model.Activity{}

	for rows.Next() {
		var activity model.Activity

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.Title,
			&activity.SportType,
			&activity.DistanceKM,
			&activity.DurationSeconds,
			&activity.AvgPaceSeconds,
			&activity.ElevationGainM,
			&activity.ActivityDate,
			&activity.Notes,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

func (r *ActivityRepository) FindByID(ctx context.Context, userID, activityID string) (*model.Activity, error) {
	query := `
		SELECT id, user_id, title, sport_type, distance_km, duration_seconds,
		       avg_pace_seconds, elevation_gain_m, activity_date, notes,
		       created_at, updated_at
		FROM activities
		WHERE user_id = $1 AND id = $2
	`

	activity := &model.Activity{}

	err := r.db.QueryRow(ctx, query, userID, activityID).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.Title,
		&activity.SportType,
		&activity.DistanceKM,
		&activity.DurationSeconds,
		&activity.AvgPaceSeconds,
		&activity.ElevationGainM,
		&activity.ActivityDate,
		&activity.Notes,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func (r *ActivityRepository) Update(ctx context.Context, activity *model.Activity) error {
	query := `
		UPDATE activities
		SET title = $1,
		    sport_type = $2,
		    distance_km = $3,
		    duration_seconds = $4,
		    avg_pace_seconds = $5,
		    elevation_gain_m = $6,
		    activity_date = $7,
		    notes = $8,
		    updated_at = NOW()
		WHERE user_id = $9 AND id = $10
		RETURNING updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		activity.Title,
		activity.SportType,
		activity.DistanceKM,
		activity.DurationSeconds,
		activity.AvgPaceSeconds,
		activity.ElevationGainM,
		activity.ActivityDate,
		activity.Notes,
		activity.UserID,
		activity.ID,
	).Scan(&activity.UpdatedAt)
}
