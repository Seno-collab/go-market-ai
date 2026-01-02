package topicapp

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/transport/response"
)

type ComboboxUseCase struct {
	Repo domain.TopicRepository
}

func NewComboboxUseCase(repo domain.TopicRepository) *ComboboxUseCase {
	return &ComboboxUseCase{
		Repo: repo,
	}
}

func (useCase *ComboboxUseCase) Execute(ctx context.Context, restaurantID int32) (*[]response.Combobox, error) {
	return useCase.Repo.GetTopicsByRestaurantCombobox(ctx, restaurantID)
}
