package db

import (
	"context"
	"go-ai/internal/identity/domain/rbac"
	sqlc "go-ai/internal/identity/infrastructure/sqlc/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RbacRepo struct {
	q *sqlc.Queries
}

func NewRbacRepo(pool *pgxpool.Pool) *RbacRepo {
	return &RbacRepo{
		q: sqlc.New(pool),
	}
}

func (ru *RbacRepo) GetUserRole(ctx context.Context, userID uuid.UUID) (rbac.UserRole, error) {
	u, err := ru.q.GetUserRole(ctx, sqlc.GetUserRoleParams{
		UserID:   userID,
		IsActive: true,
	})
	if err != nil {
		return rbac.UserRole{}, err
	}

	return rbac.UserRole{
		UserID: userID,
		Role:   *u,
	}, nil
}
