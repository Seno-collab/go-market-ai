package menuitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/pkg/utils"
)

type UpdateUseCase struct {
	Repo domain.MenuItemRepository
}

func NewUpdateUseCase(repo domain.MenuItemRepository) *UpdateUseCase {
	return &UpdateUseCase{
		Repo: repo,
	}
}

func (useCase *UpdateUseCase) Execute(ctx context.Context, req UpdateMenuItemRequest, restaurantID int32, topicID int64) error {
	url, err := utils.NewUrl(req.ImageUrl)
	if err != nil {
		return err
	}
	price, err := utils.NewMoney(req.Price)
	if err != nil {
		return err
	}
	entity, err := domain.NewMenuItem(req.Name, price, req.Type, url, req.Description, req.Sku, restaurantID, domain.TopicID(topicID))
	if err := entity.Validate(); err != nil {
		return err
	}
	return nil
}
