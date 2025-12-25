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

func (uc *GetTopicsUseCase) Execute(ctx context.Context, restaurantID int32) (*GetTopicsResponse, error) {
	rows, err := uc.Repo.GetTopics(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	items := make([]GetTopicResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapTopicResponse(row))
	}
	resp := &GetTopicsResponse{
		Items: items,
	}
	return resp, nil
}
