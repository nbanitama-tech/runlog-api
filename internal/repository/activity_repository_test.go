package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/internal/testutil/integration"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
)

func TestActivityRepository_CreateAndFindByID(t *testing.T) {
	db := integration.NewDB(t)
	defer db.Close()

	integration.Cleanup(t, db)
	defer integration.Cleanup(t, db)

	ctx := context.Background()

	userRepo := NewUserRepository(db)
	activityRepo := NewActivityRepository(db)

	user := &model.User{
		Name:         "Novandi",
		Email:        "novandi.activity@example.com",
		PasswordHash: "hashed-password",
	}

	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	activityDate := time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC)

	activity := &model.Activity{
		UserID:          user.ID,
		Title:           "Morning Easy Run",
		SportType:       "running",
		DistanceKM:      10.5,
		DurationSeconds: 3600,
		AvgPaceSeconds:  342,
		ElevationGainM:  120,
		ActivityDate:    activityDate,
		Notes:           "Easy effort",
	}

	if err := activityRepo.Create(ctx, activity); err != nil {
		t.Fatalf("create activity: %v", err)
	}

	if activity.ID == "" {
		t.Fatal("expected activity ID to be generated")
	}

	foundActivity, err := activityRepo.FindByID(ctx, user.ID, activity.ID)
	if err != nil {
		t.Fatalf("find activity: %v", err)
	}

	if foundActivity.Title != activity.Title {
		t.Fatalf("expected title %s, got %s", activity.Title, foundActivity.Title)
	}

	if foundActivity.UserID != user.ID {
		t.Fatalf("expected user id %s, got %s", user.ID, foundActivity.UserID)
	}
}

func TestActivityRepository_FindByUserID(t *testing.T) {
	db := integration.NewDB(t)
	defer db.Close()

	integration.Cleanup(t, db)
	defer integration.Cleanup(t, db)

	ctx := context.Background()

	userRepo := NewUserRepository(db)
	activityRepo := NewActivityRepository(db)

	user := &model.User{
		Name:         "Novandi",
		Email:        "novandi.list@example.com",
		PasswordHash: "hashed-password",
	}

	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	activityDate := time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC)

	activity := &model.Activity{
		UserID:          user.ID,
		Title:           "Easy Run",
		SportType:       "running",
		DistanceKM:      8.0,
		DurationSeconds: 3000,
		AvgPaceSeconds:  375,
		ElevationGainM:  50,
		ActivityDate:    activityDate,
		Notes:           "List test",
	}

	if err := activityRepo.Create(ctx, activity); err != nil {
		t.Fatalf("create activity: %v", err)
	}

	result, err := activityRepo.FindByUserID(ctx, user.ID, model.ActivityFilter{
		Page:      1,
		PageSize:  20,
		Offset:    0,
		SortBy:    "activity_date",
		SortOrder: "DESC",
	})
	if err != nil {
		t.Fatalf("find activities: %v", err)
	}

	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
	}

	if len(result.Activities) != 1 {
		t.Fatalf("expected 1 activity, got %d", len(result.Activities))
	}
}

func TestActivityRepository_Update(t *testing.T) {
	db := integration.NewDB(t)
	defer db.Close()

	integration.Cleanup(t, db)
	defer integration.Cleanup(t, db)

	ctx := context.Background()

	userRepo := NewUserRepository(db)
	activityRepo := NewActivityRepository(db)

	user := &model.User{
		Name:         "Novandi",
		Email:        "novandi.update@example.com",
		PasswordHash: "hashed-password",
	}

	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	activityDate := time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC)

	activity := &model.Activity{
		UserID:          user.ID,
		Title:           "Original Run",
		SportType:       "running",
		DistanceKM:      5.0,
		DurationSeconds: 1800,
		AvgPaceSeconds:  360,
		ElevationGainM:  20,
		ActivityDate:    activityDate,
		Notes:           "Before update",
	}

	if err := activityRepo.Create(ctx, activity); err != nil {
		t.Fatalf("create activity: %v", err)
	}

	activity.Title = "Updated Run"
	activity.DistanceKM = 10.0
	activity.DurationSeconds = 3600
	activity.AvgPaceSeconds = 360
	activity.Notes = "After update"

	if err := activityRepo.Update(ctx, activity); err != nil {
		t.Fatalf("update activity: %v", err)
	}

	updatedActivity, err := activityRepo.FindByID(ctx, user.ID, activity.ID)
	if err != nil {
		t.Fatalf("find updated activity: %v", err)
	}

	if updatedActivity.Title != "Updated Run" {
		t.Fatalf("expected Updated Run, got %s", updatedActivity.Title)
	}

	if updatedActivity.DistanceKM != 10.0 {
		t.Fatalf("expected distance 10.0, got %v", updatedActivity.DistanceKM)
	}
}

func TestActivityRepository_Delete(t *testing.T) {
	db := integration.NewDB(t)
	defer db.Close()

	integration.Cleanup(t, db)
	defer integration.Cleanup(t, db)

	ctx := context.Background()

	userRepo := NewUserRepository(db)
	activityRepo := NewActivityRepository(db)

	user := &model.User{
		Name:         "Novandi",
		Email:        "novandi.delete@example.com",
		PasswordHash: "hashed-password",
	}

	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	activityDate := time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC)

	activity := &model.Activity{
		UserID:          user.ID,
		Title:           "Delete Run",
		SportType:       "running",
		DistanceKM:      3.0,
		DurationSeconds: 1200,
		AvgPaceSeconds:  400,
		ElevationGainM:  10,
		ActivityDate:    activityDate,
		Notes:           "Delete test",
	}

	if err := activityRepo.Create(ctx, activity); err != nil {
		t.Fatalf("create activity: %v", err)
	}

	if err := activityRepo.Delete(ctx, user.ID, activity.ID); err != nil {
		t.Fatalf("delete activity: %v", err)
	}

	_, err := activityRepo.FindByID(ctx, user.ID, activity.ID)
	if err == nil {
		t.Fatal("expected error after deleting activity")
	}
}

func TestActivityRepository_FindByID_NotFound(t *testing.T) {
	db := integration.NewDB(t)
	defer db.Close()

	integration.Cleanup(t, db)
	defer integration.Cleanup(t, db)

	ctx := context.Background()
	activityRepo := NewActivityRepository(db)

	_, err := activityRepo.FindByID(
		ctx,
		"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000002",
	)

	if !errors.Is(err, pkgerrors.ErrActivityNotFound) {
		t.Fatalf("expected ErrActivityNotFound, got %v", err)
	}
}

func TestActivityRepository_Delete_NotFound(t *testing.T) {
	db := integration.NewDB(t)
	defer db.Close()

	integration.Cleanup(t, db)
	defer integration.Cleanup(t, db)

	ctx := context.Background()
	activityRepo := NewActivityRepository(db)

	err := activityRepo.Delete(
		ctx,
		"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000002",
	)

	if !errors.Is(err, pkgerrors.ErrActivityNotFound) {
		t.Fatalf("expected ErrActivityNotFound, got %v", err)
	}
}
