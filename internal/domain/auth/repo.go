package auth

import "github.com/google/uuid"

type Repository interface {
	// GetByID(id int64) (*Auth, error)
	GetByEmail(email string) (*Auth, error)
	CreateUser(u *Auth) (uuid.UUID, error)
}
