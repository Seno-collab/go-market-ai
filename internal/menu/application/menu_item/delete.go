package menuitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type DeleteUseCase struct {
	Repo domain.MenuItemRepository
}

func NewDeleteUseCase(repo domain.MenuItemRepository) *DeleteUseCase {
	return &DeleteUseCase{
		Repo: repo,
	}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id int64, restaurantID int32) error {
	return uc.Repo.DeleteMenuItem(ctx, id, restaurantID)
}
