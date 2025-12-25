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

func (uc *UpdateUseCase) Execute(ctx context.Context, req UpdateMenuItemRequest, restaurantID int32, menuItemID int64) error {
	existing, err := uc.Repo.GetMenuItemByID(ctx, menuItemID, restaurantID)
	if err != nil {
		return err
	}
	url, err := utils.NewUrl(req.ImageUrl)
	if err != nil {
		return err
	}
	price, err := utils.NewMoney(req.Price)
	if err != nil {
		return err
	}
	existing.Name = req.Name
	existing.Sku = req.Sku
	existing.Description = req.Description
	existing.Type = req.Type
	existing.ImageUrl = url
	existing.BasePrice = price
	if err := existing.Validate(); err != nil {
		return err
	}
	return uc.Repo.UpdateMenuItem(ctx, &existing)
}
