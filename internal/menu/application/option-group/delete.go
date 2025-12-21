package optiongroupapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type DeleteUseCase struct {
	Repo domain.OptionRepository
}

func NewDeleteUseCase(repo domain.OptionRepository) *DeleteUseCase {
	return &DeleteUseCase{Repo: repo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id int64, restaurantID int32) error {
	return uc.Repo.DeleteOptionGroup(ctx, id, restaurantID)
}
