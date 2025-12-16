package db

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TopicRepo struct {
	Sqlc *sqlc.Queries
	Pool *pgxpool.Pool
}

func NewTopicRepo(pool *pgxpool.Pool) *TopicRepo {
	return &TopicRepo{
		Sqlc: sqlc.New(pool),
		Pool: pool,
	}
}
func (tr *MenuRepo) GetTopics(ctx context.Context, restaurant int32) ([]domain.Topic, error) {
	rows, err := tr.Sqlc.GetTopicsByRestaurant(ctx, restaurant)
	if err != nil {
		return nil, err
	}
	topics := make([]domain.Topic, 0, len(rows))
	for _, row := range rows {
		topics = append(topics, domain.Topic{
			ID:           domain.TopicID(row.ID),
			RestaurantID: row.RestaurantID,
			Name:         row.Name,
			Slug:         *row.Slug,
			// ParentID:     (*domain.TopicID)(row.ParentID),
			SortOrder: row.SortOrder,
			IsActive:  row.IsActive,
		})
	}
	return topics, nil
}

func (tr *MenuRepo) CreateTopic(ctx context.Context, t *domain.Topic) (domain.TopicID, error) {
	row, err := tr.Sqlc.CreateTopic(ctx, sqlc.CreateTopicParams{
		RestaurantID: t.RestaurantID,
		Name:         t.Name,
		Slug:         &t.Slug,
		ParentID:     (*int64)(t.ParentID),
		SortOrder:    t.SortOrder,
	})
	if err != nil {
		return 0, err
	}
	return domain.TopicID(row.ID), err
}

func (tr *MenuRepo) DeleteTopic(ctx context.Context, id domain.TopicID) error {
	return tr.Sqlc.DeleteTopic(ctx, int64(id))
}

func (tr *MenuRepo) UpdateTopic(ctx context.Context, t *domain.Topic) error {
	if err := tr.Sqlc.UpdateTopic(ctx, sqlc.UpdateTopicParams{
		Name:         t.Name,
		ID:           int64(t.ID),
		Slug:         &t.Slug,
		ParentID:     (*int64)(t.ParentID),
		RestaurantID: t.RestaurantID,
	}); err != nil {
		return err
	}
	return nil
}
