package user

import "errors"

var (
	ErrNotFound     = errors.New("user not found")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidName  = errors.New("invalid name")
	ErrConflict     = errors.New("email already exists")
)
