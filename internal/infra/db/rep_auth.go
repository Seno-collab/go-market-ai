package db

import (
	"context"
	domain "go-ai/internal/domain/auth"
	"go-ai/internal/infra/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	q *sqlc.Queries
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		q: sqlc.New(pool),
	}
}

func (r *AuthRepo) GetByEmail(email string) (*domain.Auth, error) {
	u, err := r.q.GetUserByEmail(context.Background(), &email)
	if err != nil {
		return nil, err
	}

	return &domain.Auth{
		ID:    u.ID,
		Email: *u.Email,
	}, nil
}

func (r *AuthRepo) CreateUser(a *domain.Auth) (uuid.UUID, error) {
	id, err := r.q.CreateUser(context.Background(), sqlc.CreateUserParams{
		Email:        &a.Email,
		PasswordHash: a.Password,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
