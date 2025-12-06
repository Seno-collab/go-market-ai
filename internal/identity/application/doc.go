package app

import (
	authapp "go-ai/internal/identity/application/auth"
	"go-ai/internal/transport/response"
)

type RegisterSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	Data *authapp.RegisterSuccess `json:"data,omitempty"`
}

type GetProfileSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	Data *authapp.GetProfileResponse `json:"data,omitempty"`
}

type RefreshTokenSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	Data *authapp.RefreshTokenResponse `json:"data,omitempty"`
}

type LoginSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	Data *authapp.LoginResponse `json:"data,omitempty"`
}
