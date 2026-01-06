package db

import (
	"context"
	"go-ai/internal/identity/domain/auth"
	sqlc "go-ai/internal/identity/infrastructure/sqlc/user"
	"go-ai/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	queries *sqlc.Queries
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		queries: sqlc.New(pool),
	}
}

func (au *AuthRepo) GetByEmail(ctx context.Context, email string) (*auth.Entity, error) {
	u, err := au.queries.GetUserByEmail(ctx, &email)
	if err != nil {
		return nil, err
	}
	if u.Email == nil {
		return nil, auth.ErrInvalidEmail
	}
	em, err := utils.NewEmail(*u.Email)
	if err != nil {
		return nil, auth.ErrInvalidEmail
	}
	role := ""
	if u.RoleName != nil {
		role = *u.RoleName
	}
	imageUrl := ""
	if u.ImageUrl != nil {
		imageUrl = *u.ImageUrl
	}
	pw, _ := auth.NewPasswordFromHash(u.PasswordHash)
	return &auth.Entity{
		ID:       u.ID,
		Email:    em,
		FullName: u.FullName,
		Password: pw,
		Role:     role,
		ImageUrl: imageUrl,
		IsActive: u.IsActive,
	}, nil
}

func (au *AuthRepo) CreateUser(ctx context.Context, a *auth.Entity) (uuid.UUID, error) {
	email := a.Email.String()
	password := a.Password.String()

	id, err := au.queries.CreateUser(ctx, sqlc.CreateUserParams{
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
	u, err := au.queries.GetUserByName(ctx, name)
	if err != nil {
		return nil, err
	}
	em, err := utils.NewEmail(*u.Email)
	if err != nil {
		return nil, auth.ErrInvalidEmail
	}
	imageUrl := ""
	if u.ImageUrl != nil {
		imageUrl = *u.ImageUrl
	}
	return &auth.Entity{
		ID:       u.ID,
		Email:    em,
		FullName: u.FullName,
		Role:     *u.RoleName,
		ImageUrl: imageUrl,
		IsActive: u.IsActive,
	}, nil
}

func (au *AuthRepo) GetById(ctx context.Context, id uuid.UUID) (*auth.Entity, error) {
	u, err := au.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u.ID == uuid.Nil {
		return nil, auth.ErrUserNotFound
	}
	em, err := utils.NewEmail(*u.Email)
	imageUrl := ""
	if u.ImageUrl != nil {
		imageUrl = *u.ImageUrl
	}
	return &auth.Entity{
		ID:       u.ID,
		Email:    em,
		FullName: u.FullName,
		Role:     *u.RoleName,
		ImageUrl: imageUrl,
		IsActive: u.IsActive,
	}, nil
}

func (au *AuthRepo) ChangePassword(ctx context.Context, NewPasswordHash string, userID uuid.UUID) error {
	return au.queries.UpdatePasswordByID(ctx, sqlc.UpdatePasswordByIDParams{
		PasswordHash: NewPasswordHash,
		ID:           userID,
	})
}

func (au *AuthRepo) GetPasswordByID(ctx context.Context, id uuid.UUID) (string, error) {
	return au.queries.GetPasswordByID(ctx, id)
}

func (au *AuthRepo) UpdateProfile(ctx context.Context, u *auth.Entity) error {
	email := u.Email.String()
	password, err := au.GetPasswordByID(ctx, u.ID)
	if err != nil {
		return err
	}
	return au.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		FullName:     u.FullName,
		Email:        &email,
		PasswordHash: password,
		ImageUrl:     &u.ImageUrl,
		IsActive:     u.IsActive,
		ID:           u.ID,
	})
}
