package optionitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type DeleteUseCase struct {
	Repo domain.OptionItemRepository
}

func NewDeleteUseCase(repo domain.OptionItemRepository) *DeleteUseCase {
	return &DeleteUseCase{Repo: repo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id int64, restaurantID int32) error {
	return uc.Repo.DeleteOptionItem(ctx, id, restaurantID)
}
