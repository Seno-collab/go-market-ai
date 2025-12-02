package restaurantapp

import "go-ai/internal/domain/restaurant"

type UpdateStatusUseCase struct {
	repo restaurant.Repository
}

func NewUpdateStatusUseCase(repo restaurant.Repository) *UpdateStatusUseCase {
	return &UpdateStatusUseCase{
		repo: repo,
	}
}

func (uc *UpdateStatusUseCase) Execute(id int32, dayOfWeek int32) error {
	return nil
}
