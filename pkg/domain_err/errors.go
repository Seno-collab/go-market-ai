package domainerr

import (
	"net/http"
)

var (
	ErrInvalidEmail        = New(http.StatusBadRequest, "invalid email")
	ErrInternalServerError = New(http.StatusInternalServerError, "internal server error")
	ErrInvalidField        = New(http.StatusBadRequest, "invalid field")
)
