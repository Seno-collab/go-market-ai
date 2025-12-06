package app

import (
	restaurantapp "go-ai/internal/restaurant/application/restaurant"
	"go-ai/internal/transport/response"
)

type CreateRestaurantSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	Data *restaurantapp.CreateRestaurantResponse `json:"data,omitempty"`
}

type GetRestaurantByIDSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	Data *restaurantapp.GetRestaurantByIDResponse `json:"data,omitempty"`
}

type UpdateRestaurantSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}

type DeleteRestaurantSuccessResponseDoc struct {
	response.SuccecssBaseDoc
}
