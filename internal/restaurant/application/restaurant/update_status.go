package restaurantapp

import "go-ai/internal/restaurant/domain/restaurant"

type UpdateStatusUseCase struct {
	Repo restaurant.Repository
}

func NewUpdateStatusUseCase(repo restaurant.Repository) *UpdateStatusUseCase {
	return &UpdateStatusUseCase{
		Repo: repo,
	}
}

func (uc *UpdateStatusUseCase) Execute(id int32, dayOfWeek int32) error {
	return nil
}
