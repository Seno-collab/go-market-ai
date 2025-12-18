package app

import (
	menuitemapp "go-ai/internal/menu/application/menu_item"
	topicapp "go-ai/internal/menu/application/topic"
	"go-ai/internal/transport/response"
)

type CreateMenuItemSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type CreateTopicSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type GetTopicSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type GetMenuItemSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	Data *menuitemapp.GetMenuItemResponse `json:"data,omitempty"`
}

type UpdateMenuItemSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type DeleteMenuItemSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type GetMenuItemsByRestaurantSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	menuitemapp.GetMenuItemsResponse
}

type GetTopicsByRestaurantSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	*topicapp.GetTopicsResponse
}
