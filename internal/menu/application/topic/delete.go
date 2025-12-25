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

func (uc *DeleteUseCase) Execute(ctx context.Context, id int64, restaurantID int32) error {
	return uc.Repo.DeleteTopic(ctx, domain.TopicID(id), restaurantID)
}
