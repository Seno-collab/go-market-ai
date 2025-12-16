package topicapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type GetUseCase struct {
	Repo domain.TopicRepository
}

func NewGetUseCase(repo domain.TopicRepository) *CreateUseCase {
	return &CreateUseCase{
		Repo: repo,
	}
}

func (useCase *GetUseCase) Execute(ctx context.Context, id int64) error {
	return nil
}
