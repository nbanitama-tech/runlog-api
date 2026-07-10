package model

import "time"

// ActivityFilter represents the filter criteria for listing activities, including pagination parameters (Page, PageSize, Offset), sport type, date range (From, To), and sorting options (SortBy, SortOrder). It is used to encapsulate the filtering and sorting options when querying activities from the repository.
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

// ActivityListResult represents the result of listing activities, including a slice of Activity objects and the total count of activities. It is used to encapsulate the response from the activity repository when retrieving a list of activities based on specific filters and pagination parameters.
type ActivityListResult struct {
	Activities []Activity
	Total      int
}
