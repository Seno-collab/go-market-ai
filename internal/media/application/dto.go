package app

import (
	uploadapp "go-ai/internal/media/application/upload"
	"go-ai/internal/transport/response"
)

type UploadLogoSuccessResponseDoc struct {
	response.SuccecssBaseDoc
	Data *uploadapp.UploadLogoResponse `json:"data,omitempty"`
}
