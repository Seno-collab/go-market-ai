package variantapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type CreateUseCase struct {
	Repo domain.VariantsRepository
}

func NewCreateUseCase(repo domain.VariantsRepository) *CreateUseCase {
	return &CreateUseCase{
		Repo: repo,
	}
}

func (useCase *CreateUseCase) Execute(ctx context.Context) error {

	return nil
}
