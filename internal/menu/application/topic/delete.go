package topicapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type DeleteUseCase struct {
	Repo domain.TopicRepository
}

func NewDeleteUseCase(repo domain.TopicRepository) *DeleteUseCase {
	return &DeleteUseCase{
		Repo: repo,
	}
}

func (useCase *DeleteUseCase) Execute(ctx context.Context, id int64, restaurantID int32) error {
	return useCase.Repo.DeleteTopic(ctx, domain.TopicID(id), restaurantID)
}
