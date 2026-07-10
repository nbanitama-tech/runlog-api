package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/model"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
)

type mockActivityUseCase struct {
	activity *model.Activity
	result   *model.ActivityListResult
	err      error
}

func (m *mockActivityUseCase) Create(_ context.Context, _, _, _ string, _ float64, _, _ int, _ time.Time, _ string) (*model.Activity, error) {
	return m.activity, m.err
}

func (m *mockActivityUseCase) ListByUserID(_ context.Context, _ string, _ model.ActivityFilter) (*model.ActivityListResult, error) {
	return m.result, m.err
}

func (m *mockActivityUseCase) GetByID(_ context.Context, _, _ string) (*model.Activity, error) {
	return m.activity, m.err
}

func (m *mockActivityUseCase) Update(_ context.Context, _, _, _, _ string, _ float64, _, _ int, _ time.Time, _ string) (*model.Activity, error) {
	return m.activity, m.err
}

func (m *mockActivityUseCase) Delete(_ context.Context, _, _ string) error {
	return m.err
}

func TestActivityHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockActivityUseCase{
		activity: &model.Activity{
			ID:              "activity-123",
			UserID:          "user-123",
			Title:           "Morning Easy Run",
			SportType:       "running",
			DistanceKM:      10.5,
			DurationSeconds: 3600,
			AvgPaceSeconds:  342,
			ElevationGainM:  120,
			ActivityDate:    time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC),
			Notes:           "Easy effort",
			CreatedAt:       time.Now(),
		},
	}

	h := NewActivityHandler(mockUC)

	r := gin.New()
	r.POST("/activities", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		h.Create(c)
	})

	body := map[string]any{
		"title":            "Morning Easy Run",
		"sport_type":       "running",
		"distance_km":      10.5,
		"duration_seconds": 3600,
		"elevation_gain_m": 120,
		"activity_date":    "2026-07-09",
		"notes":            "Easy effort",
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/activities", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestActivityHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockActivityUseCase{
		result: &model.ActivityListResult{
			Activities: []model.Activity{
				{
					ID:              "activity-123",
					UserID:          "user-123",
					Title:           "Morning Easy Run",
					SportType:       "running",
					DistanceKM:      10.5,
					DurationSeconds: 3600,
					AvgPaceSeconds:  342,
					ElevationGainM:  120,
					ActivityDate:    time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC),
				},
			},
			Total: 1,
		},
	}

	h := NewActivityHandler(mockUC)

	r := gin.New()
	r.GET("/activities", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		h.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/activities?page=1&page_size=10", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestActivityHandler_Detail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockActivityUseCase{
		activity: &model.Activity{
			ID:              "activity-123",
			UserID:          "user-123",
			Title:           "Morning Easy Run",
			SportType:       "running",
			DistanceKM:      10.5,
			DurationSeconds: 3600,
			AvgPaceSeconds:  342,
			ElevationGainM:  120,
			ActivityDate:    time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC),
		},
	}

	h := NewActivityHandler(mockUC)

	r := gin.New()
	r.GET("/activities/:id", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		h.Detail(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/activities/activity-123", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestActivityHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockActivityUseCase{
		activity: &model.Activity{
			ID:              "activity-123",
			UserID:          "user-123",
			Title:           "Updated Run",
			SportType:       "running",
			DistanceKM:      12.0,
			DurationSeconds: 4200,
			AvgPaceSeconds:  350,
			ElevationGainM:  150,
			ActivityDate:    time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC),
			UpdatedAt:       time.Now(),
		},
	}

	h := NewActivityHandler(mockUC)

	r := gin.New()
	r.PUT("/activities/:id", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		h.Update(c)
	})

	body := map[string]any{
		"title":            "Updated Run",
		"sport_type":       "running",
		"distance_km":      12.0,
		"duration_seconds": 4200,
		"elevation_gain_m": 150,
		"activity_date":    "2026-07-09",
		"notes":            "Updated notes",
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/activities/activity-123", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestActivityHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockActivityUseCase{}

	h := NewActivityHandler(mockUC)

	r := gin.New()
	r.DELETE("/activities/:id", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		h.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/activities/activity-123", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestActivityHandler_Detail_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockActivityUseCase{
		err: pkgerrors.ErrActivityNotFound,
	}

	h := NewActivityHandler(mockUC)

	r := gin.New()
	r.GET("/activities/:id", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		h.Detail(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/activities/not-found", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestActivityHandler_Create_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &mockActivityUseCase{}

	h := NewActivityHandler(mockUC)

	r := gin.New()
	r.POST("/activities", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		h.Create(c)
	})

	body := map[string]any{
		"title":            "",
		"distance_km":      0,
		"duration_seconds": 0,
		"activity_date":    "bad-date",
	}

	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/activities", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", w.Code, w.Body.String())
	}
}
