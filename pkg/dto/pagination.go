package dto

// PaginationMeta represents the metadata for paginated responses. It includes fields for the current page number (Page), the number of items per page (PageSize), the total number of items (Total), and the total number of pages (TotalPages). The struct is used to provide pagination information in API responses when returning a list of activities in the RunLog API application.
type PaginationMeta struct {
	Page       int `json:"page" example:"1"`
	PageSize   int `json:"page_size" example:"20"`
	Total      int `json:"total" example:"100"`
	TotalPages int `json:"total_pages" example:"5"`
}

// PaginatedActivityResponseEnvelope represents the response envelope for a paginated list of activities. It includes a success flag, a slice of activity data, and pagination metadata. The struct is used to wrap the paginated list of activity responses in a consistent format for API responses in the RunLog API application.
type PaginatedActivityResponseEnvelope struct {
	Success bool               `json:"success" example:"true"`
	Data    []ActivityResponse `json:"data"`
	Meta    PaginationMeta     `json:"meta"`
}
