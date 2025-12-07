package auth

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (*Entity, error)
	GetByEmail(ctx context.Context, email string) (*Entity, error)
	CreateUser(ctx context.Context, u *Entity) (uuid.UUID, error)
	GetByName(ctx context.Context, name string) (*Entity, error)
	ChangePassword(ctx context.Context, password string, userID uuid.UUID) error
	GetPasswordByID(ctx context.Context, userID uuid.UUID) (string, error)
}
