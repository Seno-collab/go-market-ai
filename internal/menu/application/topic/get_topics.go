package topicapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type GetTopicsUseCase struct {
	Repo domain.TopicRepository
}

func NewGetTopicsUseCase(repo domain.TopicRepository) *GetTopicsUseCase {
	return &GetTopicsUseCase{
		Repo: repo,
	}
}

func (useCase *GetTopicsUseCase) Execute(ctx context.Context, restaurantID int32) (*GetTopicsResponse, error) {
	rows, err := useCase.Repo.GetTopics(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	items := make([]GetTopicResponse, len(rows))
	for _, row := range rows {
		items = append(items, GetTopicResponse{
			Name:         row.Name,
			RestaurantID: row.RestaurantID,
			Slug:         row.Slug,
			SortOrder:    row.SortOrder,
			ParentID:     int64(*row.ParentID),
		})
	}
	resp := &GetTopicsResponse{
		Items: items,
	}
	return resp, nil
}
