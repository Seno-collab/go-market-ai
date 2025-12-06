package rbac

import (
	"context"

	"github.com/google/uuid"
)

type UserRoleRepo interface {
	GetUserRole(ctx context.Context, userID uuid.UUID) (UserRole, error)
}
