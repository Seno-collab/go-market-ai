package optionitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/pkg/utils"
)

type UpdateUseCase struct {
	Repo domain.OptionItemRepository
}

func NewUpdateUseCase(repo domain.OptionItemRepository) *UpdateUseCase {
	return &UpdateUseCase{Repo: repo}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, id int64, req UpdateOptionItemRequest, restaurantID int32) error {
	if err := validateQuantity(req.QuantityMin, req.QuantityMax); err != nil {
		return err
	}
	existing, err := uc.Repo.GetOptionItem(ctx, id, restaurantID)
	if err != nil {
		return err
	}
	price, err := utils.NewMoney(req.PriceDelta)
	if err != nil {
		return err
	}
	existing.Name = req.Name
	existing.LinkedMenuItem = req.LinkedMenuItem
	existing.PriceDelta = price
	existing.QuantityMin = req.QuantityMin
	existing.QuantityMax = req.QuantityMax
	existing.SortOrder = req.SortOrder
	return uc.Repo.UpdateOptionItem(ctx, &existing, restaurantID)
}
