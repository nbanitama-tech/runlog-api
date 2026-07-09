package handler

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/model"
	"github.com/nbanitama-tech/runlog-api/pkg/dto"
	pkgerrors "github.com/nbanitama-tech/runlog-api/pkg/errors"
	"github.com/nbanitama-tech/runlog-api/pkg/response"
)

type ActivityUseCase interface {
	Create(ctx context.Context, userID, title, sportType string, distanceKM float64, durationSeconds, elevationGainM int, activityDate time.Time, notes string) (*model.Activity, error)
	ListByUserID(ctx context.Context, userID string, filter model.ActivityFilter) (*model.ActivityListResult, error)
	GetByID(ctx context.Context, userID, activityID string) (*model.Activity, error)
	Update(ctx context.Context, userID, activityID, title, sportType string, distanceKM float64, durationSeconds, elevationGainM int, activityDate time.Time, notes string) (*model.Activity, error)
	Delete(ctx context.Context, userID, activityID string) error
}

type ActivityHandler struct {
	activityUseCase ActivityUseCase
}

func NewActivityHandler(activityUseCase ActivityUseCase) *ActivityHandler {
	return &ActivityHandler{activityUseCase: activityUseCase}
}

// CreateActivity godoc
//
//	@Summary		Create activity
//	@Description	Create a new running activity
//	@Tags			Activities
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateActivityRequest	true	"Activity"
//	@Success		201		{object}	dto.ActivityResponseEnvelope
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Router			/activities [post]
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

// ListActivities godoc
//
//	@Summary		List activities
//	@Description	Get all activities for the authenticated user
//	@Tags			Activities
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Param			sport_type	query		string	false	"Sport type"
//	@Param			from		query		string	false	"From date YYYY-MM-DD"
//	@Param			to			query		string	false	"To date YYYY-MM-DD"
//	@Param			sort	query	string	false	"Sort field. Use -field for descending. Example: -activity_date"
//	@Success		200			{object}	dto.PaginatedActivityResponseEnvelope
//	@Success		200	{array}		dto.ActivityResponseEnvelope
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/activities [get]
func (h *ActivityHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	var query dto.ListActivityQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	filter, err := query.ToFilter()
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.activityUseCase.ListByUserID(c.Request.Context(), userID, filter)
	if err != nil {
		response.InternalServerError(c)
		return
	}

	totalPages := 0
	if result.Total > 0 {
		totalPages = (result.Total + filter.PageSize - 1) / filter.PageSize
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    dto.ToActivityResponses(result.Activities),
		"meta": dto.PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      result.Total,
			TotalPages: totalPages,
		},
	})
}

// GetActivity godoc
//
//	@Summary		Get activity
//	@Description	Get activity by ID
//	@Tags			Activities
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Activity ID"
//	@Success		200	{object}	dto.ActivityResponseEnvelope
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/activities/{id} [get]

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

// UpdateActivity godoc
//
//	@Summary		Update activity
//	@Description	Update an existing activity
//	@Tags			Activities
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Activity ID"
//	@Param			request	body		dto.UpdateActivityRequest	true	"Activity"
//	@Success		200		{object}	dto.ActivityResponseEnvelope
//	@Failure		404		{object}	dto.ErrorResponse
//	@Router			/activities/{id} [put]
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

// DeleteActivity godoc
//
//	@Summary		Delete activity
//	@Description	Delete activity by ID
//	@Tags			Activities
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"Activity ID"
//	@Success		204
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/activities/{id} [delete]
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
