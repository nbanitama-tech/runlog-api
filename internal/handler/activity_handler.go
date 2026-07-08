package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/usecase"
	"github.com/nbanitama-tech/runlog-api/pkg/dto"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
	"github.com/nbanitama-tech/runlog-api/pkg/response"
)

type ActivityHandler struct {
	activityUseCase *usecase.ActivityUseCase
}

func NewActivityHandler(activityUseCase *usecase.ActivityUseCase) *ActivityHandler {
	return &ActivityHandler{activityUseCase: activityUseCase}
}

func (h *ActivityHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req dto.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	activityDate, err := time.Parse("2006-01-02", req.ActivityDate)
	if err != nil {
		response.BadRequest(c, "activity_date must use YYYY-MM-DD format")
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
		response.InternalServerError(c)
		return
	}

	response.Created(c, dto.ToActivityResponse(*activity))
}

func (h *ActivityHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	activities, err := h.activityUseCase.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.OK(c, dto.ToActivityResponses(activities))
}

func (h *ActivityHandler) Detail(c *gin.Context) {
	userID := c.GetString("user_id")
	activityID := c.Param("id")

	activity, err := h.activityUseCase.GetByID(c.Request.Context(), userID, activityID)
	if err != nil {
		if errors.Is(err, pkgerrors.ErrActivityNotFound) {
			response.NotFound(c, "activity not found")
			return
		}
		response.InternalServerError(c)

		return
	}

	response.OK(c, dto.ToActivityResponse(*activity))
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
		response.BadRequest(c, "invalid request body")
		return
	}

	activityDate, err := time.Parse("2006-01-02", req.ActivityDate)
	if err != nil {
		response.BadRequest(c, "activity_date must use YYYY-MM-DD format")
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
		if errors.Is(err, pkgerrors.ErrActivityNotFound) {
			response.NotFound(c, "activity not found")
			return
		}
		response.InternalServerError(c)
		return
	}

	response.OK(c, dto.ToActivityResponse(*activity))
}

func (h *ActivityHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	activityID := c.Param("id")

	err := h.activityUseCase.Delete(c.Request.Context(), userID, activityID)
	if err != nil {
		if errors.Is(err, pkgerrors.ErrActivityNotFound) {
			response.NotFound(c, "activity not found")
			return
		}

		response.InternalServerError(c)
		return
	}

	response.NoContent(c)
}
