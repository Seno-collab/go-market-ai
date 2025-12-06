package restaurantapp

import (
	"context"
	"go-ai/internal/restaurant/domain/restaurant"

	"github.com/google/uuid"
)

type DeleteUseCase struct {
	Repo restaurant.Repository
}

func NewDeleteUseCase(repo restaurant.Repository) *DeleteUseCase {
	return &DeleteUseCase{
		Repo: repo,
	}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id int32, userID uuid.UUID) error {
	err := uc.Repo.SoftDelete(ctx, id, userID)
	if err != nil {
		return err
	}
	return nil
}
