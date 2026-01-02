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

func (uc *CreateUseCase) Execute(ctx context.Context, req CreateTopicRequest, restaurantID int32) error {
	entity, err := newTopicEntity(restaurantID, req.Name, req.Slug, req.ParentID, req.SortOrder)
	if err != nil {
		return err
	}
	if err := entity.Validate(); err != nil {
		return err
	}
	_, err = uc.Repo.CreateTopic(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

func newTopicEntity(restaurantID int32, name, slug string, parentValue int64, sortOrder int32) (*domain.Topic, error) {
	if parentValue == 0 {
		parentValue = 1
	}
	var parent *domain.TopicID
	if parentValue != 0 {
		val := domain.TopicID(parentValue)
		parent = &val
	}
	return domain.NewTopic(restaurantID, name, slug, parent, sortOrder)
}
