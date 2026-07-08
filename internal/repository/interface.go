package repository

import (
	"context"

	"github.com/nbanitama-tech/runlog-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type ActivityRepository interface {
	Create(ctx context.Context, activity *model.Activity) error
	FindByUserID(ctx context.Context, userID string) ([]model.Activity, error)
	FindByID(ctx context.Context, userID, activityID string) (*model.Activity, error)
	Update(ctx context.Context, activity *model.Activity) error
	Delete(ctx context.Context, userID, activityID string) error
}
