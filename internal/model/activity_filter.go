package model

import "time"

type ActivityFilter struct {
	Page      int
	PageSize  int
	Offset    int
	SportType string
	From      *time.Time
	To        *time.Time
	SortBy    string
	SortOrder string
}

type ActivityListResult struct {
	Activities []Activity
	Total      int
}
