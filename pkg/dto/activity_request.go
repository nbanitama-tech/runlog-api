package dto

type ListActivityQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	SportType string `form:"sport_type"`
	From      string `form:"from"`
	To        string `form:"to"`
	Sort      string `form:"sort"`
}
