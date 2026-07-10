package dto

// ListActivityQuery represents the query parameters for listing activities. It includes fields for pagination (Page and PageSize), filtering by sport type (SportType), date range (From and To), and sorting options (Sort). The struct is used to capture and validate query parameters from API requests when retrieving a list of activities.
type ListActivityQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	SportType string `form:"sport_type"`
	From      string `form:"from"`
	To        string `form:"to"`
	Sort      string `form:"sort"`
}
