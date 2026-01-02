package topicapp

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/transport/response"
	"math"
)

type GetTopicsUseCase struct {
	Repo domain.TopicRepository
}

func NewGetTopicsUseCase(repo domain.TopicRepository) *GetTopicsUseCase {
	return &GetTopicsUseCase{
		Repo: repo,
	}
}

func (uc *GetTopicsUseCase) Execute(ctx context.Context, restaurantID int32, name string, page, limit int32) (*GetTopicsResponse, error) {
	page, limit, offset := response.ApplyDefaultPaginated(&page, &limit)
	rows, total, err := uc.Repo.GetTopics(ctx, name, restaurantID, limit, offset)
	if err != nil {
		return nil, err
	}
	items := make([]GetTopicResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapTopicResponse(row))
	}
	totalPages := int64(math.Ceil(float64(total) / float64(limit)))
	return &GetTopicsResponse{
		PaginatedResponse: response.PaginatedResponse[[]GetTopicResponse]{
			Items:      items,
			Page:       page,
			Limit:      limit,
			TotalItems: int64(len(items)),
			TotalPages: totalPages,
		},
	}, nil
}
