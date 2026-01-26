package menuitemapp

import (
	"go-ai/internal/menu/domain"
	"go-ai/internal/transport/response"
	"time"
)

type CreateMenuItemRequest struct {
	Name        string              `json:"name"`
	Sku         string              `json:"sku"`
	Description string              `json:"description"`
	Type        domain.MenuItemType `json:"type"`
	ImageUrl    string              `json:"image_url"`
	Price       int64               `json:"price"`
	TopicID     domain.TopicID      `json:"topic_id"`
}

type GetMenuItemResponse struct {
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	Sku         string              `json:"sku"`
	Description string              `json:"description"`
	Type        domain.MenuItemType `json:"type"`
	ImageUrl    string              `json:"image_url"`
	Price       int64               `json:"price"`
	IsActive    bool                `json:"is_active"`
	UpdatedAt   time.Time           `json:"updated_at"`
	CreateAt    time.Time           `json:"created_at"`
}

type OptionItemDetail struct {
	ID             int64  `json:"id"`
	OptionGroupID  int64  `json:"option_group_id"`
	Name           string `json:"name"`
	LinkedMenuItem *int64 `json:"linked_menu_item,omitempty"`
	PriceDelta     int64  `json:"price_delta"`
	QuantityMin    int32  `json:"quantity_min"`
	QuantityMax    int32  `json:"quantity_max"`
	SortOrder      int32  `json:"sort_order"`
	IsActive       bool   `json:"is_active"`
}

type OptionGroupDetail struct {
	ID         int64              `json:"id"`
	Name       string             `json:"name"`
	MinSelect  int32              `json:"min_select"`
	MaxSelect  int32              `json:"max_select"`
	IsRequired bool               `json:"is_required"`
	SortOrder  int32              `json:"sort_order"`
	Items      []OptionItemDetail `json:"items"`
}

type MenuItemDetailResponse struct {
	GetMenuItemResponse
	OptionGroups []OptionGroupDetail `json:"option_groups"`
}

type UpdateMenuItemRequest struct {
	Name        string              `json:"name"`
	Sku         string              `json:"sku"`
	Description string              `json:"description"`
	Type        domain.MenuItemType `json:"type"`
	ImageUrl    string              `json:"image_url"`
	Price       int64               `json:"price"`
}

type GetMenuItemsRequest struct {
	Page     *int32  `json:"page,omitempty"`
	Limit    *int32  `json:"limit,omitempty"`
	Filter   *string `json:"filter,omitempty"`
	Category *string `json:"category,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type GetMenuItemsResponse struct {
	response.PaginatedResponse[[]GetMenuItemResponse]
}

type UpdateMenuItemStatusRequest struct {
	IsActive bool `json:"is_active"`
}
