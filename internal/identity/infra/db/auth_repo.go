package db

import (
	"context"
	"go-ai/internal/identity/domain/auth"
	sqlc "go-ai/internal/identity/infra/sqlc/user"
	"go-ai/pkg/utils"

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

func (au *AuthRepo) GetByEmail(ctx context.Context, email string) (*auth.Entity, error) {
	u, err := au.q.GetUserByEmail(ctx, &email)
	if err != nil {
		return nil, err
	}
	em, err := utils.NewEmail(*u.Email)
	if err != nil {
		return nil, auth.ErrInvalidEmail
	}
	role := ""
	if u.RoleName != nil {
		role = *u.RoleName
	}
	pw, _ := auth.NewPasswordFromHash(u.PasswordHash)
	return &auth.Entity{
		ID:       u.ID,
		Email:    em,
		FullName: u.FullName,
		Password: pw,
		Role:     role,
		IsActive: u.IsActive,
	}, nil
}

func (au *AuthRepo) CreateUser(ctx context.Context, a *auth.Entity) (uuid.UUID, error) {
	email := a.Email.String()
	password := a.Password.String()

	id, err := au.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        &email,
		PasswordHash: password,
		FullName:     a.FullName,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
func (au *AuthRepo) GetByName(ctx context.Context, name string) (*auth.Entity, error) {
	u, err := au.q.GetUserByName(ctx, name)
	if err != nil {
		return nil, err
	}
	em, err := utils.NewEmail(*u.Email)
	return &auth.Entity{
		ID:       u.ID,
		Email:    em,
		FullName: u.FullName,
		Role:     *u.RoleName,
		IsActive: u.IsActive,
	}, nil
}

func (au *AuthRepo) GetById(ctx context.Context, id uuid.UUID) (*auth.Entity, error) {
	u, err := au.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	em, err := utils.NewEmail(*u.Email)
	return &auth.Entity{
		ID:       u.ID,
		Email:    em,
		FullName: u.FullName,
		Role:     *u.RoleName,
		IsActive: u.IsActive,
	}, nil
}
