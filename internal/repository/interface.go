package repository

import (
	"context"

	"github.com/nbanitama-tech/runlog-api/internal/model"
)

// UserRepository defines the interface for user-related database operations. It provides methods for creating a new user and finding a user by their email address. Implementations of this interface should handle the underlying database interactions to manage user data in the RunLog API application.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

// ActivityRepository defines the interface for activity-related database operations. It provides methods for creating, listing, retrieving, updating, and deleting user activities. Implementations of this interface should handle the underlying database interactions to manage activity data in the RunLog API application.
type ActivityRepository interface {
	Create(ctx context.Context, activity *model.Activity) error
	FindByUserID(ctx context.Context, userID string, filter model.ActivityFilter) (*model.ActivityListResult, error)
	FindByID(ctx context.Context, userID, activityID string) (*model.Activity, error)
	Update(ctx context.Context, activity *model.Activity) error
	Delete(ctx context.Context, userID, activityID string) error
}
