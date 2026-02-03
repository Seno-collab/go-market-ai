package authapp

import (
	"go-ai/internal/transport/response"
)

type RegisterSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type GetProfileSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *GetProfileResponse `json:"data,omitempty"`
}

type UpdateProfileSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *GetProfileResponse `json:"data,omitempty"`
}

type RefreshTokenSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *RefreshTokenResponse `json:"data,omitempty"`
}

type LoginSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *LoginResponse `json:"data,omitempty"`
}

type ChangePasswordSuccessResponseDoc struct {
	response.SuccessBaseDoc
}

type LogoutSuccessResponseDoc struct {
	response.SuccessBaseDoc
}
