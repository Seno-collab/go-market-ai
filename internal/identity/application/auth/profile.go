package authapp

import (
	"context"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infrastructure/cache"
	domainerr "go-ai/pkg/domain_err"

	"time"

	"github.com/google/uuid"
)

type GetProfileUseCase struct {
	Repo  auth.Repository
	Cache *cache.AuthCache
}

func NewGetProfileUseCase(repo auth.Repository, cache *cache.AuthCache) *GetProfileUseCase {
	return &GetProfileUseCase{
		Repo:  repo,
		Cache: cache,
	}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, userID uuid.UUID) (*GetProfileResponse, error) {
	keyAuth := fmt.Sprintf("profile_%s", userID.String())
	cacheData, err := uc.Cache.GetAuthCache(ctx, keyAuth)
	if err != nil {
		return nil, domainerr.ErrInternalServerError
	}
	profile := &GetProfileResponse{}
	if cacheData == nil {
		record, err := uc.Repo.GetById(ctx, userID)
		if err != nil {
			return nil, auth.ErrUserNotFound
		}
		profile = &GetProfileResponse{
			Email:    record.Email.String(),
			FullName: record.FullName,
			Role:     record.Role,
			IsActive: record.IsActive,
		}
		authData := &cache.AuthData{
			Email:    record.Email.String(),
			FullName: record.FullName,
			Role:     record.Role,
			IsActive: record.IsActive,
		}
		uc.Cache.SetAuthCache(keyAuth, authData, time.Duration(60*int(time.Minute)))
		return profile, nil
	}
	profile = &GetProfileResponse{
		Email:    cacheData.Email,
		FullName: cacheData.FullName,
		Role:     cacheData.Role,
		IsActive: cacheData.IsActive,
	}
	return profile, nil
}
