package restaurantapp

import (
	"context"
	"go-ai/internal/restaurant/domain/restaurant"
	"go-ai/pkg/utils"

	"github.com/google/uuid"
)

type UpdateRestaurantUseCase struct {
	Repo restaurant.Repository
}

func NewUpdateRestaurantUseCase(repo restaurant.Repository) *UpdateRestaurantUseCase {
	return &UpdateRestaurantUseCase{
		Repo: repo,
	}
}

func (uc *UpdateRestaurantUseCase) Execute(ctx context.Context, req CreateRestaurantRequest, userID uuid.UUID, id int32) error {
	record, err := uc.Repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	if record == nil {
		return restaurant.ErrRestaurantNotExists
	}
	// if record.CreatedBy != userID {
	// 	return response.ErrForbidden
	// }
	email, err := utils.NewEmail(req.Email)
	if err != nil {
		return err
	}
	phone, err := restaurant.NewPhone(req.PhoneNumber)
	if err != nil {
		return err
	}
	logoUrl, err := restaurant.NewUrl(req.LogoUrl)
	if err != nil {
		return err
	}

	bannerUrl, err := restaurant.NewUrl(req.BannerUrl)
	if err != nil {
		return err
	}

	websiteUrl, err := restaurant.NewUrl(req.WebsiteUrl)
	if err != nil {
		return err
	}
	hours := make([]restaurant.Hours, 0, len(req.Hours))
	for _, hour := range req.Hours {
		hours = append(hours, restaurant.Hours{
			Day:       hour.Day,
			OpenTime:  hour.OpenTime,
			CloseTime: hour.CloseTime,
		})
	}
	if err := restaurant.ValidateHours(hours); err != nil {
		return err
	}
	updated := &restaurant.Entity{
		ID:          record.ID,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Category:    req.Category,
		City:        req.City,
		District:    req.District,
		LogoUrl:     logoUrl,
		BannerUrl:   bannerUrl,
		PhoneNumber: phone,
		WebsiteUrl:  websiteUrl,
		Email:       email,
		UpdateBy:    userID,
		Hours:       hours,
	}

	err = uc.Repo.Update(ctx, updated, id)
	if err != nil {
		return err
	}
	return nil
}
