package uploadapp

import (
	"go-ai/pkg/response"
)

type UploadLogoSuccessResponseDoc struct {
	response.SuccessBaseDoc
	Data *UploadLogoResponse `json:"data,omitempty"`
}
