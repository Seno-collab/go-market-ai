package optiongroupapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type CreateUseCase struct {
	OptionRepo   domain.OptionRepository
	MenuItemRepo domain.MenuItemRepository
}

func NewCreateUseCase(optionRepo domain.OptionRepository, menuItemRepo domain.MenuItemRepository) *CreateUseCase {
	return &CreateUseCase{
		OptionRepo:   optionRepo,
		MenuItemRepo: menuItemRepo,
	}
}

func (uc *CreateUseCase) Execute(ctx context.Context, req CreateOptionGroupRequest, restaurantID int32) (int64, error) {
	if _, err := uc.MenuItemRepo.GetMenuItemByID(ctx, req.MenuItemID, restaurantID); err != nil {
		return 0, err
	}
	group, err := domain.NewOptionGroup(restaurantID, req.Name, req.MinSelect, req.MaxSelect, req.IsRequired, req.SortOrder)
	if err != nil {
		return 0, err
	}
	id, err := uc.OptionRepo.CreateOptionGroup(ctx, group)
	if err != nil {
		return 0, err
	}
	if err := uc.OptionRepo.AttachOptionGroupToItem(ctx, req.MenuItemID, id, req.SortOrder); err != nil {
		return 0, err
	}
	return id, nil
}
