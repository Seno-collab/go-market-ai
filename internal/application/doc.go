package app

import (
	authapp "go-ai/internal/application/auth"
	restaurantapp "go-ai/internal/application/restaurant"
	uploadapp "go-ai/internal/application/upload"
	"go-ai/internal/transport/http/response"
)

type ErrorResponseDoc struct {
	ResponseCode string                `json:"response_code,omitempty"`
	Message      string                `json:"message"`
	Error        *response.ErrorDetail `json:"error,omitempty"`
}

type SuccecssResponseBaseDoc struct {
	Message      string `json:"message"`
	ResponseCode string `json:"response_code,omitempty"`
}

type RegisterSuccessResponseDoc struct {
	SuccecssResponseBaseDoc
	Data *authapp.RegisterSuccess `json:"data,omitempty"`
}

type GetProfileSuccessResponseDoc struct {
	SuccecssResponseBaseDoc
	Data *authapp.GetProfileResponse `json:"data,omitempty"`
}

type RefreshTokenSuccessResponseDoc struct {
	SuccecssResponseBaseDoc
	Data *authapp.RefreshTokenResponse `json:"data,omitempty"`
}

type LoginSuccessResponseDoc struct {
	SuccecssResponseBaseDoc
	Data *authapp.LoginResponse `json:"data,omitempty"`
}

type UploadLogoSuccessResponseDoc struct {
	SuccecssResponseBaseDoc
	Data *uploadapp.UploadLogoResponse `json:"data,omitempty"`
}

type CreateRestaurantSuccessResponseDoc struct {
	SuccecssResponseBaseDoc
}

type GetRestaurantByIDSuccessResponseDoc struct {
	SuccecssResponseBaseDoc
	Data *restaurantapp.GetRestaurantByIDResponse `json:"data,omitempty"`
}
