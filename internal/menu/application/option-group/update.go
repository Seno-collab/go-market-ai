package optiongroupapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type UpdateUseCase struct {
	Repo domain.OptionRepository
}

func NewUpdateUseCase(repo domain.OptionRepository) *UpdateUseCase {
	return &UpdateUseCase{Repo: repo}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, id int64, req UpdateOptionGroupRequest, restaurantID int32) error {
	group, err := domain.NewOptionGroup(restaurantID, req.Name, req.MinSelect, req.MaxSelect, req.IsRequired, req.SortOrder)
	if err != nil {
		return err
	}
	group.ID = id
	return uc.Repo.UpdateOptionGroup(ctx, group)
}
