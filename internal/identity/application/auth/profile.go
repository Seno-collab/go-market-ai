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

type GetProfileResponse struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
	ImageUrl string `json:"image_url"`
}

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
	sessionKey := fmt.Sprintf("session_%s", userID.String())
	cacheData, err := uc.Cache.GetAuthCache(ctx, sessionKey)
	if err != nil {
		return nil, domainerr.ErrInternalServerError
	}
	profile := &GetProfileResponse{}
	if cacheData == nil {
		record, err := uc.Repo.GetById(ctx, userID)
		if err != nil {
			return nil, auth.ErrUserNotFound
		}
		authData := &cache.UserCache{
			UserID:   userID,
			Email:    record.Email.String(),
			FullName: record.FullName,
			Role:     record.Role,
			IsActive: record.IsActive,
			ImageUrl: record.ImageUrl,
		}
		uc.Cache.SetAuthCache(ctx, sessionKey, authData, time.Duration(60*int(time.Minute)))
	}
	profile = &GetProfileResponse{
		Email:    cacheData.Email,
		FullName: cacheData.FullName,
		Role:     cacheData.Role,
		IsActive: cacheData.IsActive,
		ImageUrl: cacheData.ImageUrl,
	}
	return profile, nil
}
