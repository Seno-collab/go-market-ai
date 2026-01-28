package authapp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infrastructure/cache"
	domainerr "go-ai/pkg/domain_err"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UpdateProfileUseCase struct {
	Repo  auth.Repository
	Cache *cache.AuthCache
}

func NewUpdateProfileUseCase(repo auth.Repository, cache *cache.AuthCache) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		Repo:  repo,
		Cache: cache,
	}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*GetProfileResponse, error) {
	user, err := uc.Repo.GetById(ctx, userID)
	if err != nil {
		if ae, ok := err.(domainerr.AppError); ok {
			return nil, ae
		}
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrUserNotFound
		}
		return nil, domainerr.ErrInternalServerError
	}

	if req.FullName != "" && req.FullName != user.FullName {
		if err := user.UpdateFullName(req.FullName); err != nil {
			return nil, err
		}
	}

	if req.Email != "" && req.Email != user.Email.String() {
		if err := user.UpdateEmail(req.Email); err != nil {
			return nil, err
		}

		existing, err := uc.Repo.GetByEmail(ctx, req.Email)
		if err == nil && existing != nil && existing.ID != user.ID {
			return nil, auth.ErrEmailAlreadyExists
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrInternalServerError
		}
	}

	if req.ImageUrl != "" && req.ImageUrl != user.ImageUrl {
		if err := user.UpdateImageUrl(req.ImageUrl); err != nil {
			return nil, err
		}
	}

	if err := uc.Repo.UpdateProfile(ctx, user); err != nil {
		return nil, domainerr.ErrInternalServerError
	}

	profile := &GetProfileResponse{
		Email:    user.Email.String(),
		FullName: user.FullName,
		Role:     user.Role,
		IsActive: user.IsActive,
		ImageUrl: user.ImageUrl,
	}

	sessionKey := fmt.Sprintf("session_%s", userID.String())
	uc.Cache.SetAuthCache(ctx, sessionKey, &cache.UserCache{
		UserID:   userID,
		Email:    user.Email.String(),
		FullName: user.FullName,
		Role:     user.Role,
		IsActive: user.IsActive,
		ImageUrl: user.ImageUrl,
	}, time.Duration(60*int(time.Minute)))

	return profile, nil
}
