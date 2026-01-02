package domainerr

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidEmail        = New(http.StatusBadRequest, "Invalid email")
	ErrInternalServerError = New(http.StatusInternalServerError, "Internal server error")
	ErrInvalidField        = New(http.StatusBadRequest, "Invalid field")
	ErrInvalidUrl          = New(http.StatusBadRequest, "Invalid url")
	ErrInvalidPrice        = errors.New("Menu: price must be >= 0")
)
