package app

import (
	menuitemapp "go-ai/internal/menu/application/menu_item"
	optiongroupapp "go-ai/internal/menu/application/option-group"
	optionitemapp "go-ai/internal/menu/application/option_item"
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

type UpdateTopicSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type DeleteTopicSuccessResponseDoc struct {
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

type CreateVariantSuccessRespons struct {
	response.SuccecssBaseDoc
}

type CreateOptionGroupSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	*optiongroupapp.CreateOptionGroupResponse
}

type GetOptionGroupSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	*optiongroupapp.GetOptionGroupResponse
}

type GetOptionGroupsSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	*optiongroupapp.GetOptionGroupsResponse
}

type UpdateOptionGroupSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type DeleteOptionGroupSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type CreateOptionItemSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	*optionitemapp.CreateOptionItemResponse
}

type GetOptionItemSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	*optionitemapp.GetOptionItemResponse
}

type GetOptionItemsSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	*optionitemapp.GetOptionItemsResponse
}

type UpdateOptionItemSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type DeleteOptionItemSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}
