package restaurant

import "context"

type Repository interface {
	Create(ctx context.Context, r *Entity) (int32, error)
	GetById(ctx context.Context, id int32) (*Entity, error)
	GetByName(ctx context.Context, name string) (*Entity, error)
}
