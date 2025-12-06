package authapp

import (
	"context"
	"database/sql"
	"errors"

	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infra/cache"
	"go-ai/internal/platform/security"
	"go-ai/internal/transport/response"
	"strings"

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
	if !strings.Contains(request.Email, "@") {
		return uuid.Nil, response.ErrInvalidEmail
	}
	if request.FullName == "" {
		return uuid.Nil, response.ErrInvalidName
	}
	if request.Password == "" {
		return uuid.Nil, response.ErrInvalidPassword
	}
	// unique email
	record, err := s.Repo.GetByEmail(ctx, request.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, err
		}
	}
	if record != nil {
		return uuid.Nil, response.ErrUserAlreadyExists
	}
	record, err = s.Repo.GetByName(ctx, request.FullName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, err
		}
	}
	if record != nil {
		return uuid.Nil, response.ErrNameAlreadyExists
	}
	hasedPassword, err := security.HashPassword(request.Password)
	if err != nil {
		return uuid.Nil, response.ErrInternalServerError
	}
	return s.Repo.CreateUser(ctx, &auth.Entity{
		FullName: request.FullName,
		Email:    request.Email,
		Password: hasedPassword,
	})
}
