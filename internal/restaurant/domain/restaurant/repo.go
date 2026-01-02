package restaurant

import (
	"context"
	"go-ai/internal/transport/response"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, r *Entity, userID uuid.UUID) (int32, error)
	GetById(ctx context.Context, id int32) (*Entity, error)
	GetByName(ctx context.Context, name string) (*Entity, error)
	Update(ctx context.Context, r *Entity, id int32) error
	SoftDelete(ctx context.Context, id int32, userID uuid.UUID) error
	GetRestaurantItemCombobox(ctx context.Context, userID uuid.UUID) (*[]response.Combobox, error)
}
