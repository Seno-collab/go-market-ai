package topicapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type GetUseCase struct {
	Repo domain.TopicRepository
}

func NewGetUseCase(repo domain.TopicRepository) *GetUseCase {
	return &GetUseCase{
		Repo: repo,
	}
}

func (useCase *GetUseCase) Execute(ctx context.Context, id int64, restaurantID int32) (*GetTopicResponse, error) {
	row, err := useCase.Repo.GetTopic(ctx, domain.TopicID(id), restaurantID)
	if err != nil {
		return nil, err
	}
	resp := mapTopicResponse(row)
	return &resp, nil
}

func mapTopicResponse(row domain.Topic) GetTopicResponse {
	var parent *int64
	if row.ParentID != nil {
		val := int64(*row.ParentID)
		parent = &val
	}
	return GetTopicResponse{
		ID:           int64(row.ID),
		RestaurantID: row.RestaurantID,
		Name:         row.Name,
		Slug:         row.Slug,
		ParentID:     parent,
		SortOrder:    row.SortOrder,
	}
}
