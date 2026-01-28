package authapp

import (
	"context"
	"database/sql"
	"errors"

	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infrastructure/cache"
	"go-ai/internal/platform/security"
	domainerr "go-ai/pkg/domain_err"
	"go-ai/pkg/utils"

	"github.com/google/uuid"
)

type RegisterUseCase struct {
	Repo  auth.Repository
	Cache *cache.AuthCache
}

func NewRegisterUseCase(repo auth.Repository, cache *cache.AuthCache) *RegisterUseCase {
	return &RegisterUseCase{
		Repo:  repo,
		Cache: cache,
	}
}

func (s *RegisterUseCase) Execute(ctx context.Context, request RegisterRequest) (uuid.UUID, error) {

	email, err := utils.NewEmail(request.Email)
	if err != nil {
		return uuid.Nil, err
	}

	// unique email
	record, err := s.Repo.GetByEmail(ctx, request.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, err
		}
	}
	if record != nil {
		return uuid.Nil, auth.ErrEmailAlreadyExists
	}
	record, err = s.Repo.GetByName(ctx, request.FullName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, err
		}
	}
	if record != nil {
		return uuid.Nil, auth.ErrNameAlreadyExists
	}

	rawPassword, err := auth.NewPassword(request.Password)
	if err != nil {
		return uuid.Nil, err
	}

	hashedPassword, err := security.HashPassword(rawPassword.String())
	if err != nil {
		return uuid.Nil, domainerr.ErrInternalServerError
	}
	pw, _ := auth.NewPasswordFromHash(hashedPassword)
	return s.Repo.CreateUser(ctx, &auth.Entity{
		FullName: request.FullName,
		Email:    email,
		Password: pw,
	})
}
