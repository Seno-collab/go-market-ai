package app

import (
	menuitemapp "go-ai/internal/menu/application/menu_item"
	optiongroupapp "go-ai/internal/menu/application/option-group"
	optionitemapp "go-ai/internal/menu/application/option_item"
	topicapp "go-ai/internal/menu/application/topic"
	"go-ai/internal/transport/response"
)

type CreateMenuItemSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type CreateTopicSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type GetTopicSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type UpdateTopicSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type DeleteTopicSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type GetMenuItemSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *menuitemapp.GetMenuItemResponse `json:"data,omitempty"`
}

type UpdateMenuItemSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type DeleteMenuItemSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type GetMenuItemsByRestaurantResponseDoc struct {
	response.PaginatedResponseDoc
	Items []menuitemapp.GetMenuItemResponse `json:"items"`
}

type GetMenuItemsByRestaurantSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data GetMenuItemsByRestaurantResponseDoc `json:"data"`
}

type GetTopicsByRestaurantResponseDoc struct {
	response.PaginatedResponseDoc
	Items []topicapp.GetTopicResponse `json:"items"`
}

type GetTopicsByRestaurantSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data GetTopicsByRestaurantResponseDoc `json:"data"`
}

type CreateVariantSuccessRespons struct {
	response.SuccessBaseDoc
}

type CreateOptionGroupSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *optiongroupapp.CreateOptionGroupResponse `json:"data"`
}

type GetOptionGroupSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *optiongroupapp.GetOptionGroupResponse `json:"data"`
}

type GetOptionGroupsSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *optiongroupapp.GetOptionGroupsResponse `json:"data,omitempty"`
}

type UpdateOptionGroupSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type DeleteOptionGroupSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type CreateOptionItemSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *optionitemapp.CreateOptionItemResponse `json:"data,omitempty"`
}

type GetOptionItemSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *optionitemapp.GetOptionItemResponse `json:"data,omitempty"`
}

type GetOptionItemsResponseDoc struct {
	response.PaginatedResponseDoc
	Items GetOptionItemSuccessResponseDoc `json:"items"`
}

type GetOptionItemsSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data GetOptionItemsResponseDoc `json:"data"`
}

type UpdateOptionItemSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type DeleteOptionItemSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type GetTopicsByRestaurantComboboxSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *[]response.Combobox `json:"data"`
}

type UpdateMenuItemStatusSuccessResponseDoc struct {
	response.SuccessBaseDoc
}
