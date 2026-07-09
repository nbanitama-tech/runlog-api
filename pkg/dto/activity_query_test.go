package dto

import "testing"

func TestListActivityQuery_ToFilter_Defaults(t *testing.T) {
	query := ListActivityQuery{}

	filter, err := query.ToFilter()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if filter.Page != 1 {
		t.Fatalf("expected page 1, got %d", filter.Page)
	}

	if filter.PageSize != 20 {
		t.Fatalf("expected page_size 20, got %d", filter.PageSize)
	}

	if filter.SortBy != "activity_date" {
		t.Fatalf("expected sort_by activity_date, got %s", filter.SortBy)
	}

	if filter.SortOrder != "DESC" {
		t.Fatalf("expected sort_order DESC, got %s", filter.SortOrder)
	}
}

func TestListActivityQuery_ToFilter_WithSortDescending(t *testing.T) {
	query := ListActivityQuery{
		Sort: "-distance_km",
	}

	filter, err := query.ToFilter()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if filter.SortBy != "distance_km" {
		t.Fatalf("expected sort_by distance_km, got %s", filter.SortBy)
	}

	if filter.SortOrder != "DESC" {
		t.Fatalf("expected sort_order DESC, got %s", filter.SortOrder)
	}
}

func TestListActivityQuery_ToFilter_InvalidSort(t *testing.T) {
	query := ListActivityQuery{
		Sort: "invalid_field",
	}

	_, err := query.ToFilter()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListActivityQuery_ToFilter_InvalidDate(t *testing.T) {
	query := ListActivityQuery{
		From: "09-07-2026",
	}

	_, err := query.ToFilter()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
