package auth

import (
	"strings"

	"github.com/google/uuid"
)

type Auth struct {
	ID       uuid.UUID
	FullName string
	Email    string
	Password string
	RoleId   int
}

func (u *Auth) Validate() error {
	if !strings.Contains(u.Email, "@") {
		return ErrInvalidEmail
	}
	if strings.TrimSpace(u.FullName) == "" {
		return ErrInvalidName
	}
	if strings.TrimSpace(u.Password) == "" {
		return ErrInvalidPassword
	}
	return nil
}

func (u *Auth) ValidateLogin() error {
	if strings.Contains(u.Email, "@") {
		return ErrInvalidEmail
	}
	if strings.TrimSpace(u.Password) == "" {
		return ErrInvalidPassword
	}
	return nil
}
