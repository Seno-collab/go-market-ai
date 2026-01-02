package app

import (
	authapp "go-ai/internal/identity/application/auth"
	"go-ai/internal/transport/response"
)

type RegisterSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *authapp.RegisterSuccess `json:"data,omitempty"`
}

type GetProfileSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *authapp.GetProfileResponse `json:"data,omitempty"`
}

type RefreshTokenSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *authapp.RefreshTokenResponse `json:"data,omitempty"`
}

type LoginSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *authapp.LoginResponse `json:"data,omitempty"`
}

type ChangePasswordSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type LogoutSuccessResponseDoc struct {
	response.SuccessBaseDoc
}
