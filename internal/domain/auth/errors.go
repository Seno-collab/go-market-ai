package auth

import "errors"

var (
	ErrNotFound        = errors.New("user not found")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidName     = errors.New("invalid name")
	ErrConflict        = errors.New("email already exists")
	ErrUserExists      = errors.New("user already exists")
)
