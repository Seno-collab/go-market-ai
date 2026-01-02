package app

import (
	restaurantapp "go-ai/internal/restaurant/application/restaurant"
	"go-ai/internal/transport/response"
)

type CreateRestaurantSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *restaurantapp.CreateRestaurantResponse `json:"data,omitempty"`
}

type GetRestaurantByIDSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *restaurantapp.GetRestaurantByIDResponse `json:"data,omitempty"`
}

type UpdateRestaurantSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type DeleteRestaurantSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type GetComboboxRestaurantSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *[]response.Combobox `json:"data,omitempty"`
}
