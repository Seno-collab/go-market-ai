package restaurantapp

import (
	"context"
	"errors"
	"go-ai/internal/restaurant/domain/restaurant"
	"go-ai/internal/transport/response"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CreateRestaurantUseCase struct {
	Repo restaurant.Repository
}

func NewCreateRestaurantUseCase(repo restaurant.Repository) *CreateRestaurantUseCase {
	return &CreateRestaurantUseCase{
		Repo: repo,
	}
}

func (uc *CreateRestaurantUseCase) Execute(ctx context.Context, request CreateRestaurantRequest, userID uuid.UUID) (int32, error) {
	if !strings.Contains(request.Email, "@") {
		return 0, response.ErrInvalidEmail
	}
	if request.Name == "" {
		return 0, response.ErrInvalidName
	}
	if request.Address == "" {
		return 0, response.ErrInvalidAddress
	}
	if request.BannerUrl == "" {
		return 0, response.ErrInvalidBanner
	}
	if request.LogoUrl == "" {
		return 0, response.ErrInvalidLogo
	}
	if request.BannerUrl == "" {
		return 0, response.ErrInvalidBanner
	}
	if request.PhoneNumber == "" {
		return 0, response.ErrInvalidPhoneNumber
	}

	record, err := uc.Repo.GetByName(ctx, request.Name)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return 0, err
		}
	}
	if record != nil {
		return 0, response.ErrRestaurantNameExitis
	}
	hours := make([]restaurant.Hours, 0, len(request.Hours))
	for _, hour := range request.Hours {
		hours = append(hours, restaurant.Hours{
			Day:       hour.Day,
			OpenTime:  hour.OpenTime,
			CloseTime: hour.CloseTime,
		})
	}
	id, err := uc.Repo.Create(ctx, &restaurant.Entity{
		Email:       request.Email,
		Name:        request.Name,
		Description: request.Description,
		Category:    request.Category,
		WebsiteUrl:  request.WebsiteUrl,
		LogoUrl:     request.LogoUrl,
		BannerUrl:   request.BannerUrl,
		PhoneNumber: request.PhoneNumber,
		Address:     request.PhoneNumber,
		City:        request.City,
		District:    request.District,
		UserID:      userID,
		Hours:       hours,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}
