package restaurantapp

import (
	"context"
	"go-ai/internal/domain/restaurant"
)

type DeleteUseCase struct {
	repo restaurant.Repository
}

func NewDeleteUseCase(repo restaurant.Repository) *DeleteUseCase {
	return &DeleteUseCase{
		repo: repo,
	}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id int32) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
