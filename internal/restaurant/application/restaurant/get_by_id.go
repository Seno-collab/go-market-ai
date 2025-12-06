package restaurantapp

import (
	"context"
	"go-ai/internal/restaurant/domain/restaurant"
	"go-ai/internal/transport/response"
)

type GetByIDUseCase struct {
	Repo restaurant.Repository
}

func NewGetByIDUseCase(repo restaurant.Repository) *GetByIDUseCase {
	return &GetByIDUseCase{
		Repo: repo,
	}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, id int32) (*GetRestaurantByIDResponse, error) {
	if id == 0 {
		return nil, response.ErrInvalidField
	}
	record, err := uc.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	// key := fmt.Sprintf("profile_%s", record.UserID.String())
	// profile, err := uc.Cache.GetAuthCache(key)
	// if err != nil {
	// 	return nil, err
	// }
	hours := make([]RestaurantHoursBase, 0, len(record.Hours))
	for _, hour := range record.Hours {
		hours = append(hours, RestaurantHoursBase{
			Day:       hour.Day,
			OpenTime:  hour.OpenTime,
			CloseTime: hour.CloseTime,
		})
	}
	return &GetRestaurantByIDResponse{
		RestaurantBase: RestaurantBase{
			Name:        record.Name,
			Description: record.Description,
			Address:     record.Address,
			Category:    record.Category,
			City:        record.City,
			District:    record.District,
			LogoUrl:     record.LogoUrl,
			BannerUrl:   record.BannerUrl,
			PhoneNumber: record.PhoneNumber,
			WebsiteUrl:  record.WebsiteUrl,
			Email:       record.Email,
		},
		Hours:    hours,
		IsActive: len(hours) > 0,
		// UserName: profile.FullName,
	}, nil
}
