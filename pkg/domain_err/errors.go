package domainerr

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidEmail        = New(http.StatusBadRequest, "invalid email")
	ErrInternalServerError = New(http.StatusInternalServerError, "internal server error")
	ErrInvalidField        = New(http.StatusBadRequest, "invalid field")
	ErrInvalidUrl          = New(http.StatusBadRequest, "invalid url")
	ErrInvalidPrice        = errors.New("menu: price must be >= 0")
)
