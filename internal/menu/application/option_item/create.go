package optionitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/pkg/utils"
)

type CreateUseCase struct {
	Repo      domain.OptionItemRepository
	GroupRepo domain.OptionRepository
}

func NewCreateUseCase(repo domain.OptionItemRepository, groupRepo domain.OptionRepository) *CreateUseCase {
	return &CreateUseCase{
		Repo:      repo,
		GroupRepo: groupRepo,
	}
}

func (uc *CreateUseCase) Execute(ctx context.Context, req CreateOptionItemRequest, restaurantID int32) (int64, error) {
	if _, err := uc.GroupRepo.GetOptionGroup(ctx, req.OptionGroupID, restaurantID); err != nil {
		return 0, err
	}
	if err := validateQuantity(req.QuantityMin, req.QuantityMax); err != nil {
		return 0, err
	}
	price, err := utils.NewMoney(req.PriceDelta)
	if err != nil {
		return 0, err
	}
	item, err := domain.NewOptionItem(req.Name, price)
	if err != nil {
		return 0, err
	}
	item.OptionGroupID = req.OptionGroupID
	item.LinkedMenuItem = req.LinkedMenuItem
	item.QuantityMin = req.QuantityMin
	item.QuantityMax = req.QuantityMax
	item.SortOrder = req.SortOrder
	return uc.Repo.CreateOptionItem(ctx, item)
}

func validateQuantity(min int32, max *int32) error {
	if max == nil {
		return nil
	}
	if *max < min {
		return domain.ErrOptionGroupInvalidMinMax
	}
	return nil
}
