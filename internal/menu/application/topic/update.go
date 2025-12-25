package topicapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type UpdateUseCase struct {
	Repo domain.TopicRepository
}

func NewUpdateUseCase(repo domain.TopicRepository) *UpdateUseCase {
	return &UpdateUseCase{
		Repo: repo,
	}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, id int64, req UpdateTopicRequest, restaurantID int32) error {
	entity, err := newTopicEntity(restaurantID, req.Name, req.Slug, req.ParentID, req.SortOrder)
	if err != nil {
		return err
	}
	entity.ID = domain.TopicID(id)
	if err := entity.Validate(); err != nil {
		return err
	}
	return uc.Repo.UpdateTopic(ctx, entity)
}
