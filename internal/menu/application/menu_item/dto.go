package menuitemapp

import "go-ai/internal/menu/domain"

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
	Name        string              `json:"name"`
	Sku         string              `json:"sku"`
	Description string              `json:"description"`
	Type        domain.MenuItemType `json:"type"`
	ImageUrl    string              `json:"image_url"`
	Price       int64               `json:"price"`
}

type UpdateMenuItemRequest struct {
	Name        string              `json:"name"`
	Sku         string              `json:"sku"`
	Description string              `json:"description"`
	Type        domain.MenuItemType `json:"type"`
	ImageUrl    string              `json:"image_url"`
	Price       int64               `json:"price"`
}

type GetMenuItemsResponse struct {
	Items []GetMenuItemResponse `json:"items"`
}
