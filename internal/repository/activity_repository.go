package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nbanitama-tech/runlog-api/internal/model"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
)

type activityRepository struct {
	db *pgxpool.Pool
}

func NewActivityRepository(db *pgxpool.Pool) *activityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(ctx context.Context, activity *model.Activity) error {
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

func (r *activityRepository) FindByUserID(ctx context.Context, userID string, filter model.ActivityFilter) (*model.ActivityListResult, error) {
	baseWhere := `WHERE user_id = $1`
	args := []any{userID}
	argPos := 2

	if filter.SportType != "" {
		baseWhere += ` AND sport_type = $` + strconv.Itoa(argPos)
		args = append(args, filter.SportType)
		argPos++
	}

	if filter.From != nil {
		baseWhere += ` AND activity_date >= $` + strconv.Itoa(argPos)
		args = append(args, *filter.From)
		argPos++
	}

	if filter.To != nil {
		baseWhere += ` AND activity_date <= $` + strconv.Itoa(argPos)
		args = append(args, *filter.To)
		argPos++
	}

	countQuery := `SELECT COUNT(*) FROM activities ` + baseWhere

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	orderBy := "activity_date"
	if filter.SortBy != "" {
		orderBy = filter.SortBy
	}

	orderDirection := "DESC"
	if filter.SortOrder == "ASC" {
		orderDirection = "ASC"
	}

	query := `
	SELECT id, user_id, title, sport_type, distance_km, duration_seconds,
	       avg_pace_seconds, elevation_gain_m, activity_date, notes,
	       created_at, updated_at
	FROM activities
` + baseWhere + `
	ORDER BY ` + orderBy + ` ` + orderDirection + `, created_at DESC
	LIMIT $` + strconv.Itoa(argPos) + ` OFFSET $` + strconv.Itoa(argPos+1)

	args = append(args, filter.PageSize, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []model.Activity{}

	for rows.Next() {
		var activity model.Activity

		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &model.ActivityListResult{
		Activities: activities,
		Total:      total,
	}, nil
}

func (r *activityRepository) FindByID(ctx context.Context, userID, activityID string) (*model.Activity, error) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrors.ErrActivityNotFound
		}
		return nil, err
	}

	return activity, nil
}

func (r *activityRepository) Update(ctx context.Context, activity *model.Activity) error {
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

	err := r.db.QueryRow(
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

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pkgerrors.ErrActivityNotFound
		}

	}
	return nil
}

func (r *activityRepository) Delete(ctx context.Context, userID, activityID string) error {
	query := `
		DELETE FROM activities
		WHERE user_id = $1 AND id = $2
	`

	result, err := r.db.Exec(ctx, query, userID, activityID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pkgerrors.ErrActivityNotFound
	}

	return nil
}
