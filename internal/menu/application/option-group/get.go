package optiongroupapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type GetUseCase struct {
	Repo domain.OptionRepository
}

func NewGetUseCase(repo domain.OptionRepository) *GetUseCase {
	return &GetUseCase{Repo: repo}
}

func (uc *GetUseCase) Execute(ctx context.Context, id int64, restaurantID int32) (*GetOptionGroupResponse, error) {
	group, err := uc.Repo.GetOptionGroup(ctx, id, restaurantID)
	if err != nil {
		return nil, err
	}
	resp := mapOptionGroup(group)
	return &resp, nil
}

type GetByMenuItemUseCase struct {
	Repo         domain.OptionRepository
	MenuItemRepo domain.MenuItemRepository
}

func NewGetByMenuItemUseCase(repo domain.OptionRepository, menuItemRepo domain.MenuItemRepository) *GetByMenuItemUseCase {
	return &GetByMenuItemUseCase{
		Repo:         repo,
		MenuItemRepo: menuItemRepo,
	}
}

func (uc *GetByMenuItemUseCase) Execute(ctx context.Context, menuItemID int64, restaurantID int32) (*GetOptionGroupsResponse, error) {
	if _, err := uc.MenuItemRepo.GetMenuItemByID(ctx, menuItemID, restaurantID); err != nil {
		return nil, err
	}
	groups, err := uc.Repo.GetOptionGroups(ctx, menuItemID, restaurantID)
	if err != nil {
		return nil, err
	}
	resp := GetOptionGroupsResponse{
		Groups: make([]GetOptionGroupResponse, 0, len(groups)),
	}
	for _, g := range groups {
		resp.Groups = append(resp.Groups, mapOptionGroup(g))
	}
	return &resp, nil
}

func mapOptionGroup(group domain.OptionGroup) GetOptionGroupResponse {
	return GetOptionGroupResponse{
		ID:         group.ID,
		Name:       group.Name,
		MinSelect:  group.MinSelect,
		MaxSelect:  group.MaxSelect,
		IsRequired: group.IsRequired,
		SortOrder:  group.SortOrder,
	}
}
