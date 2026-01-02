package optionitemapp

import "go-ai/internal/transport/response"

type CreateOptionItemRequest struct {
	OptionGroupID  int64  `json:"option_group_id"`
	Name           string `json:"name"`
	LinkedMenuItem *int64 `json:"linked_menu_item"`
	PriceDelta     int64  `json:"price_delta"`
	QuantityMin    int32  `json:"quantity_min"`
	QuantityMax    int32  `json:"quantity_max"`
	SortOrder      int32  `json:"sort_order"`
}

type CreateOptionItemResponse struct {
	ID int64 `json:"id"`
}

type UpdateOptionItemRequest struct {
	Name           string `json:"name"`
	LinkedMenuItem *int64 `json:"linked_menu_item"`
	PriceDelta     int64  `json:"price_delta"`
	QuantityMin    int32  `json:"quantity_min"`
	QuantityMax    int32  `json:"quantity_max"`
	SortOrder      int32  `json:"sort_order"`
}

type GetOptionItemResponse struct {
	ID             int64  `json:"id"`
	OptionGroupID  int64  `json:"option_group_id"`
	Name           string `json:"name"`
	LinkedMenuItem *int64 `json:"linked_menu_item"`
	PriceDelta     int64  `json:"price_delta"`
	QuantityMin    int32  `json:"quantity_min"`
	QuantityMax    int32  `json:"quantity_max"`
	SortOrder      int32  `json:"sort_order"`
	IsActive       bool   `json:"is_active"`
}

type GetOptionItemsRequest struct {
	Page  *int32 `json:"page,omitempty"`
	Limit *int32 `json:"limit,omitempty"`
}

type GetOptionItemsResponse struct {
	response.PaginatedResponse[[]GetOptionItemResponse]
}
