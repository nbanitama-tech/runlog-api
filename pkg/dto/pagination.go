package dto

type PaginationMeta struct {
	Page       int `json:"page" example:"1"`
	PageSize   int `json:"page_size" example:"20"`
	Total      int `json:"total" example:"100"`
	TotalPages int `json:"total_pages" example:"5"`
}

type PaginatedActivityResponseEnvelope struct {
	Success bool               `json:"success" example:"true"`
	Data    []ActivityResponse `json:"data"`
	Meta    PaginationMeta     `json:"meta"`
}
