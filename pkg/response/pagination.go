package response

import "math"

type PaginatedResponse[T any] struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
	Items      T     `json:"items,omitempty"`
}

type PaginatedResponseDoc struct {
	Page       int32 `json:"page" example:"1"`
	Limit      int32 `json:"limit" example:"10"`
	TotalItems int64 `json:"total_items" example:"42"`
	TotalPages int64 `json:"total_pages" example:"5"`
}

var (
	pageDefault  int32 = 1
	limitDefault int32 = 10
)

func (p *PaginatedResponse[T]) ApplyDefault() {
	if p.Page <= 0 {
		p.Page = pageDefault
	}
	if p.Limit <= 0 {
		p.Limit = limitDefault
	}
}

func ApplyDefaultPaginated(page, limit *int32) (int32, int32, int32) {
	pageValue := pageDefault
	limitValue := limitDefault
	if page != nil && *page > 0 {
		pageValue = *page
	}
	if limit != nil && *limit > 0 && *limit <= 100 {
		limitValue = *limit
	}
	offset := (pageValue - 1) * limitValue
	return int32(pageValue), int32(limitValue), int32(offset)
}

func CalculateTotalPages(total, limit int64) int64 {
	return int64(math.Ceil(float64(total) / float64(limit)))

}
