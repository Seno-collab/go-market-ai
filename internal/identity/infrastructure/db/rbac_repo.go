package db

import (
	"context"
	"go-ai/internal/identity/domain/rbac"
	sqlc "go-ai/internal/identity/infrastructure/sqlc/user"
	"go-ai/pkg/metrics"
	"time"

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
	start := time.Now()
	u, err := ru.q.GetUserRole(ctx, sqlc.GetUserRoleParams{
		UserID:   userID,
		IsActive: true,
	})
	metrics.RecordDBQuery("select", "user_roles", time.Since(start).Seconds(), err)
	if err != nil {
		return rbac.UserRole{}, err
	}

	return rbac.UserRole{
		UserID: userID,
		Role:   *u,
	}, nil
}
