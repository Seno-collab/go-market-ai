package menuitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/pkg/utils"
)

type CreateUseCase struct {
	Repo      domain.MenuItemRepository
	TopicRepo domain.TopicRepository
}

func NewCreateUseCase(repo domain.MenuItemRepository, topicRepo domain.TopicRepository) *CreateUseCase {
	return &CreateUseCase{
		Repo:      repo,
		TopicRepo: topicRepo,
	}
}

func (uc *CreateUseCase) Execute(ctx context.Context, req CreateMenuItemRequest, restaurantID int32) error {
	url, err := utils.NewUrl(req.ImageUrl)
	if err != nil {
		return err
	}
	price, err := utils.NewMoney(req.Price)
	if err != nil {
		return err
	}
	if !req.Type.Valid() {
		return domain.ErrRecordUpdateFailed
	}
	// Validate topic exists for this restaurant
	if req.TopicID != 0 {
		_, err := uc.TopicRepo.GetTopic(ctx, req.TopicID, restaurantID)
		if err != nil {
			return domain.ErrTopicNotFound
		}
	}
	entity, err := domain.NewMenuItem(domain.NewMenuItemParams{
		Name:         req.Name,
		Price:        price,
		Type:         req.Type,
		ImageURL:     url,
		Description:  req.Description,
		SKU:          req.Sku,
		RestaurantID: restaurantID,
		TopicID:      req.TopicID,
	})
	if err != nil {
		return err
	}
	if err := entity.Validate(); err != nil {
		return err
	}
	_, err = uc.Repo.CreateMenuItem(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}
