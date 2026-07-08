package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/usecase"
)

type ActivityHandler struct {
	activityUseCase *usecase.ActivityUseCase
}

func NewActivityHandler(activityUseCase *usecase.ActivityUseCase) *ActivityHandler {
	return &ActivityHandler{activityUseCase: activityUseCase}
}

type CreateActivityRequest struct {
	Title           string  `json:"title" binding:"required"`
	SportType       string  `json:"sport_type"`
	DistanceKM      float64 `json:"distance_km" binding:"required,gt=0"`
	DurationSeconds int     `json:"duration_seconds" binding:"required,gt=0"`
	ElevationGainM  int     `json:"elevation_gain_m"`
	ActivityDate    string  `json:"activity_date" binding:"required"`
	Notes           string  `json:"notes"`
}

func (h *ActivityHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activityDate, err := time.Parse("2006-01-02", req.ActivityDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activity_date must use YYYY-MM-DD format"})
		return
	}

	activity, err := h.activityUseCase.Create(
		c.Request.Context(),
		userID,
		req.Title,
		req.SportType,
		req.DistanceKM,
		req.DurationSeconds,
		req.ElevationGainM,
		activityDate,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create activity"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":               activity.ID,
		"title":            activity.Title,
		"sport_type":       activity.SportType,
		"distance_km":      activity.DistanceKM,
		"duration_seconds": activity.DurationSeconds,
		"avg_pace_seconds": activity.AvgPaceSeconds,
		"elevation_gain_m": activity.ElevationGainM,
		"activity_date":    activity.ActivityDate.Format("2006-01-02"),
		"notes":            activity.Notes,
		"created_at":       activity.CreatedAt,
	})
}

func (h *ActivityHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	activities, err := h.activityUseCase.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get activities"})
		return
	}

	response := []gin.H{}

	for _, activity := range activities {
		response = append(response, gin.H{
			"id":               activity.ID,
			"title":            activity.Title,
			"sport_type":       activity.SportType,
			"distance_km":      activity.DistanceKM,
			"duration_seconds": activity.DurationSeconds,
			"avg_pace_seconds": activity.AvgPaceSeconds,
			"elevation_gain_m": activity.ElevationGainM,
			"activity_date":    activity.ActivityDate.Format("2006-01-02"),
			"notes":            activity.Notes,
			"created_at":       activity.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *ActivityHandler) Detail(c *gin.Context) {
	userID := c.GetString("user_id")
	activityID := c.Param("id")

	activity, err := h.activityUseCase.GetByID(c.Request.Context(), userID, activityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               activity.ID,
		"title":            activity.Title,
		"sport_type":       activity.SportType,
		"distance_km":      activity.DistanceKM,
		"duration_seconds": activity.DurationSeconds,
		"avg_pace_seconds": activity.AvgPaceSeconds,
		"elevation_gain_m": activity.ElevationGainM,
		"activity_date":    activity.ActivityDate.Format("2006-01-02"),
		"notes":            activity.Notes,
		"created_at":       activity.CreatedAt,
		"updated_at":       activity.UpdatedAt,
	})
}

type UpdateActivityRequest struct {
	Title           string  `json:"title" binding:"required"`
	SportType       string  `json:"sport_type"`
	DistanceKM      float64 `json:"distance_km" binding:"required,gt=0"`
	DurationSeconds int     `json:"duration_seconds" binding:"required,gt=0"`
	ElevationGainM  int     `json:"elevation_gain_m"`
	ActivityDate    string  `json:"activity_date" binding:"required"`
	Notes           string  `json:"notes"`
}

func (h *ActivityHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	activityID := c.Param("id")

	var req UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	activityDate, err := time.Parse("2006-01-02", req.ActivityDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "activity_date must use YYYY-MM-DD format"})
		return
	}

	activity, err := h.activityUseCase.Update(
		c.Request.Context(),
		userID,
		activityID,
		req.Title,
		req.SportType,
		req.DistanceKM,
		req.DurationSeconds,
		req.ElevationGainM,
		activityDate,
		req.Notes,
	)
	if err != nil {
		c.JSON(404, gin.H{"error": "activity not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":               activity.ID,
		"title":            activity.Title,
		"sport_type":       activity.SportType,
		"distance_km":      activity.DistanceKM,
		"duration_seconds": activity.DurationSeconds,
		"avg_pace_seconds": activity.AvgPaceSeconds,
		"elevation_gain_m": activity.ElevationGainM,
		"activity_date":    activity.ActivityDate.Format("2006-01-02"),
		"notes":            activity.Notes,
		"updated_at":       activity.UpdatedAt,
	})
}
