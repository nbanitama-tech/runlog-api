package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/nbanitama-tech/runlog-api/internal/model"
)

type mockActivityRepository struct {
	createdActivity *model.Activity
	updatedActivity *model.Activity
	deletedID       string
	activities      []model.Activity
	activity        *model.Activity
	err             error
}

func (m *mockActivityRepository) Create(ctx context.Context, activity *model.Activity) error {
	m.createdActivity = activity
	return m.err
}

func (m *mockActivityRepository) FindByUserID(ctx context.Context, userID string, filter model.ActivityFilter) (*model.ActivityListResult, error) {
	return &model.ActivityListResult{
		Activities: m.activities,
		Total:      len(m.activities),
	}, m.err
}

func (m *mockActivityRepository) FindByID(ctx context.Context, userID, activityID string) (*model.Activity, error) {
	return m.activity, m.err
}

func (m *mockActivityRepository) Update(ctx context.Context, activity *model.Activity) error {
	m.updatedActivity = activity
	return m.err
}

func (m *mockActivityRepository) Delete(ctx context.Context, userID, activityID string) error {
	m.deletedID = activityID
	return m.err
}

func TestActivityUseCase_Create(t *testing.T) {
	repo := &mockActivityRepository{}
	uc := NewActivityUseCase(repo)

	activityDate := time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC)

	activity, err := uc.Create(
		context.Background(),
		"user-123",
		"Morning Easy Run",
		"running",
		10.0,
		3600,
		120,
		activityDate,
		"Easy effort",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if activity.AvgPaceSeconds != 360 {
		t.Fatalf("expected avg pace 360, got %d", activity.AvgPaceSeconds)
	}

	if repo.createdActivity == nil {
		t.Fatal("expected activity to be created")
	}
}

func TestActivityUseCase_ListByUserID(t *testing.T) {
	repo := &mockActivityRepository{
		activities: []model.Activity{
			{ID: "activity-1", UserID: "user-123", Title: "Easy Run"},
			{ID: "activity-2", UserID: "user-123", Title: "Long Run"},
		},
	}

	uc := NewActivityUseCase(repo)

	activities, err := uc.ListByUserID(context.Background(), "user-123", model.ActivityFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(activities.Activities) != 2 {
		t.Fatalf("expected 2 activities, got %d", len(activities.Activities))
	}
}

func TestActivityUseCase_GetByID(t *testing.T) {
	repo := &mockActivityRepository{
		activity: &model.Activity{
			ID:     "activity-1",
			UserID: "user-123",
			Title:  "Easy Run",
		},
	}

	uc := NewActivityUseCase(repo)

	activity, err := uc.GetByID(context.Background(), "user-123", "activity-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if activity.ID != "activity-1" {
		t.Fatalf("expected activity-1, got %s", activity.ID)
	}
}

func TestActivityUseCase_Update(t *testing.T) {
	repo := &mockActivityRepository{}
	uc := NewActivityUseCase(repo)

	activityDate := time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC)

	activity, err := uc.Update(
		context.Background(),
		"user-123",
		"activity-1",
		"Updated Run",
		"running",
		12.0,
		4200,
		150,
		activityDate,
		"Updated notes",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if activity.AvgPaceSeconds != 350 {
		t.Fatalf("expected avg pace 350, got %d", activity.AvgPaceSeconds)
	}

	if repo.updatedActivity == nil {
		t.Fatal("expected activity to be updated")
	}
}

func TestActivityUseCase_Delete(t *testing.T) {
	repo := &mockActivityRepository{}
	uc := NewActivityUseCase(repo)

	err := uc.Delete(context.Background(), "user-123", "activity-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.deletedID != "activity-1" {
		t.Fatalf("expected deleted activity-1, got %s", repo.deletedID)
	}
}

func TestActivityUseCase_Create_AvgPaceCalculation(t *testing.T) {
	tests := []struct {
		name                string
		distanceKM          float64
		durationSeconds     int
		expectedAvgPaceSecs int
	}{
		{
			name:                "10K in 60 minutes",
			distanceKM:          10.0,
			durationSeconds:     3600,
			expectedAvgPaceSecs: 360,
		},
		{
			name:                "5K in 25 minutes",
			distanceKM:          5.0,
			durationSeconds:     1500,
			expectedAvgPaceSecs: 300,
		},
		{
			name:                "Half marathon in 2 hours",
			distanceKM:          21.1,
			durationSeconds:     7200,
			expectedAvgPaceSecs: 341,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockActivityRepository{}
			uc := NewActivityUseCase(repo)

			activity, err := uc.Create(
				context.Background(),
				"user-123",
				"Test Run",
				"running",
				tt.distanceKM,
				tt.durationSeconds,
				0,
				time.Now(),
				"",
			)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if activity.AvgPaceSeconds != tt.expectedAvgPaceSecs {
				t.Fatalf(
					"expected avg pace %d, got %d",
					tt.expectedAvgPaceSecs,
					activity.AvgPaceSeconds,
				)
			}
		})
	}
}
