package menuitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type UpdateStatusUseCase struct {
	Repo domain.MenuItemRepository
}

func NewUpdateStatusUseCase(repo domain.MenuItemRepository) *UpdateStatusUseCase {
	return &UpdateStatusUseCase{
		Repo: repo,
	}
}

func (useCase *UpdateStatusUseCase) Execute(ctx context.Context, restaurantID int32, id int64, req UpdateMenuItemStatusRequest) error {
	return useCase.Repo.UpdateStatusMenuItem(ctx, restaurantID, id, req.IsActive)
}
