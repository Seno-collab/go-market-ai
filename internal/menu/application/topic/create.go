package topicapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type CreateUseCase struct {
	Repo domain.TopicRepository
}

func NewCreateUseCase(repo domain.TopicRepository) *CreateUseCase {
	return &CreateUseCase{
		Repo: repo,
	}
}

func (useCase *CreateUseCase) Execute(ctx context.Context, req CreateTopicRequest, restaurantID int32) error {
	entity, err := domain.NewTopic(restaurantID, req.Name, req.Slug, nil, 0)
	if err != nil {
		return err
	}
	if err := entity.Validate(); err != nil {
		return err
	}
	_, err = useCase.Repo.CreateTopic(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}
