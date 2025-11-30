package restaurantapp

import (
	"context"
	"errors"
	"go-ai/internal/domain/restaurant"
	"go-ai/internal/transport/http/status"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CreateRestaurantUseCase struct {
	repo restaurant.Repository
}

func NewCreateRestaurantUseCase(repo restaurant.Repository) *CreateRestaurantUseCase {
	return &CreateRestaurantUseCase{
		repo: repo,
	}
}

func (uc *CreateRestaurantUseCase) Execute(ctx context.Context, request CreateRestaurantRequest, userID uuid.UUID) (int32, error) {
	if !strings.Contains(request.Email, "@") {
		return 0, status.ErrInvalidEmail
	}
	if request.Name == "" {
		return 0, status.ErrInvalidName
	}
	if request.Address == "" {
		return 0, restaurant.ErrInvalidAddress
	}
	if request.BannerUrl == "" {
		return 0, restaurant.ErrInvalidBanner
	}
	if request.LogoUrl == "" {
		return 0, restaurant.ErrInvalidLogo
	}
	if request.BannerUrl == "" {
		return 0, restaurant.ErrInvalidBanner
	}
	if request.PhoneNumber == "" {
		return 0, restaurant.ErrInvalidPhoneNumber
	}
	record, err := uc.repo.GetByName(ctx, request.Name)
	if err != nil || !errors.Is(err, pgx.ErrNoRows) || record != nil {
		return 0, err
	}
	id, err := uc.repo.Create(ctx, &restaurant.Entity{
		Email:       &request.Email,
		Name:        request.Name,
		WebsiteUrl:  &request.WebsiteUrl,
		LogoUrl:     &request.LogoUrl,
		BannerUrl:   &request.BannerUrl,
		PhoneNumber: &request.PhoneNumber,
		Address:     &request.PhoneNumber,
		City:        &request.City,
		District:    &request.District,
		UserID:      userID,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}
