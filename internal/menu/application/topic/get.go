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
	resp := &GetTopicResponse{
		RestaurantID: row.RestaurantID,
		Name:         row.Name,
		Slug:         row.Slug,
		SortOrder:    row.SortOrder,
		// ParentID:     int64(*row.ParentID),
	}
	return resp, nil
}
