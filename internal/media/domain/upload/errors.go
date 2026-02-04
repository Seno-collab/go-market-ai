package upload

import (
	domainerr "go-ai/pkg/domain_err"
	"net/http"
)

var (
	ErrFileTooLarge    = domainerr.New(http.StatusBadRequest, "File size exceeds 5MB")
	ErrInvalidFileType = domainerr.New(http.StatusBadRequest, "Only PNG, JPEG, WebP allowed")
	ErrFileRequired    = domainerr.New(http.StatusBadRequest, "File is required")
	ErrFileReadFailed  = domainerr.New(http.StatusBadRequest, "Unable to read file")
	ErrFileOpenFailed  = domainerr.New(http.StatusBadRequest, "Unable to open file")
	ErrUploadFailed    = domainerr.New(http.StatusInternalServerError, "Upload to storage failed")
)
