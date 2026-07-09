package dto

import (
	"errors"
	"strings"
	"time"

	"github.com/nbanitama-tech/runlog-api/internal/model"
)

func (q ListActivityQuery) ToFilter() (model.ActivityFilter, error) {
	page := q.Page
	if page <= 0 {
		page = 1
	}

	pageSize := q.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var from *time.Time
	if q.From != "" {
		parsedFrom, err := time.Parse("2006-01-02", q.From)
		if err != nil {
			return model.ActivityFilter{}, errors.New("from must use YYYY-MM-DD format")
		}
		from = &parsedFrom
	}

	var to *time.Time
	if q.To != "" {
		parsedTo, err := time.Parse("2006-01-02", q.To)
		if err != nil {
			return model.ActivityFilter{}, errors.New("to must use YYYY-MM-DD format")
		}
		to = &parsedTo
	}

	sortBy := "activity_date"
	sortOrder := "DESC"

	allowedSorts := map[string]bool{
		"activity_date":    true,
		"distance_km":      true,
		"duration_seconds": true,
		"created_at":       true,
	}

	if q.Sort != "" {
		sortValue := q.Sort

		if strings.HasPrefix(sortValue, "-") {
			sortOrder = "DESC"
			sortValue = strings.TrimPrefix(sortValue, "-")
		} else {
			sortOrder = "ASC"
		}

		if !allowedSorts[sortValue] {
			return model.ActivityFilter{}, errors.New("invalid sort field")
		}

		sortBy = sortValue
	}

	return model.ActivityFilter{
		Page:      page,
		PageSize:  pageSize,
		Offset:    (page - 1) * pageSize,
		SportType: q.SportType,
		From:      from,
		To:        to,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}
