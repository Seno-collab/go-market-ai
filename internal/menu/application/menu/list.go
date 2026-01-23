package menuapp

import (
	"context"
	"go-ai/internal/menu/domain"
	"strconv"
)

const (
	defaultMenuPageSize int32 = 20
	maxMenuPageSize     int32 = 50
)

type ListMenusRequest struct {
	RestaurantID int32
	Type         domain.MenuType
	Topics       []string
	Limit        int32
	Cursor       *int64
}

type MenuResponse struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Type      domain.MenuType `json:"type"`
	ImageURL  string          `json:"imageUrl"`
	BasePrice string          `json:"basePrice"`
	Topics    []string        `json:"topics,omitempty"`
}

type ListMenusResponse struct {
	Items      []MenuResponse `json:"items"`
	NextCursor *int64         `json:"nextCursor,omitempty"`
}

type ListMenusUseCase struct {
	Repo domain.MenuRepository
}

func NewListMenusUseCase(repo domain.MenuRepository) *ListMenusUseCase {
	return &ListMenusUseCase{Repo: repo}
}

func normalizeMenuLimit(limit int32) int32 {
	if limit <= 0 {
		return defaultMenuPageSize
	}
	if limit > maxMenuPageSize {
		return maxMenuPageSize
	}
	return limit
}

func (uc *ListMenusUseCase) Execute(ctx context.Context, req ListMenusRequest) (*ListMenusResponse, error) {
	limit := normalizeMenuLimit(req.Limit)
	menus, nextCursor, err := uc.Repo.ListMenus(ctx, domain.ListMenusParams{
		RestaurantID: req.RestaurantID,
		Type:         req.Type,
		Topics:       req.Topics,
		Limit:        limit,
		Cursor:       req.Cursor,
	})
	if err != nil {
		return nil, err
	}

	items := make([]MenuResponse, 0, len(menus))
	for _, m := range menus {
		items = append(items, MenuResponse{
			ID:        m.ID,
			Name:      m.Name,
			Type:      m.Type,
			ImageURL:  m.ImageURL.String(),
			BasePrice: strconv.FormatInt(int64(m.BasePrice), 10),
			Topics:    m.Topics,
		})
	}
	return &ListMenusResponse{
		Items:      items,
		NextCursor: nextCursor,
	}, nil
}
