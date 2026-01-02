package app

import (
	uploadapp "go-ai/internal/media/application/upload"
	"go-ai/internal/transport/response"
)

type UploadLogoSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *uploadapp.UploadLogoResponse `json:"data,omitempty"`
}
