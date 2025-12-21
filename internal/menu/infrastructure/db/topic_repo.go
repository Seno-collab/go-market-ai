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
func (tr *TopicRepo) GetTopics(ctx context.Context, restaurantID int32) ([]domain.Topic, error) {
	rows, err := tr.Sqlc.GetTopicsByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	topics := make([]domain.Topic, 0, len(rows))
	for _, row := range rows {
		slug := ""
		if row.Slug != nil {
			slug = *row.Slug
		}
		var parent *domain.TopicID
		if row.ParentID != nil {
			val := domain.TopicID(*row.ParentID)
			parent = &val
		}
		topics = append(topics, domain.Topic{
			ID:           domain.TopicID(row.ID),
			RestaurantID: row.RestaurantID,
			Name:         row.Name,
			Slug:         slug,
			ParentID:     parent,
			SortOrder:    row.SortOrder,
			IsActive:     row.IsActive,
		})
	}
	return topics, nil
}

func (tr *TopicRepo) GetTopic(ctx context.Context, id domain.TopicID, restaurantID int32) (domain.Topic, error) {
	row, err := tr.Sqlc.GetTopic(ctx, sqlc.GetTopicParams{
		ID:           int64(id),
		RestaurantID: restaurantID,
	})
	if err != nil {
		return domain.Topic{}, err
	}
	slug := ""
	if row.Slug != nil {
		slug = *row.Slug
	}
	var parent *domain.TopicID
	if row.ParentID != nil {
		val := domain.TopicID(*row.ParentID)
		parent = &val
	}
	return domain.Topic{
		ID:           domain.TopicID(row.ID),
		RestaurantID: row.RestaurantID,
		Name:         row.Name,
		Slug:         slug,
		SortOrder:    row.SortOrder,
		IsActive:     row.IsActive,
		ParentID:     parent,
	}, nil
}

func (tr *TopicRepo) CreateTopic(ctx context.Context, t *domain.Topic) (domain.TopicID, error) {
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

func (tr *TopicRepo) UpdateTopic(ctx context.Context, t *domain.Topic) error {
	if err := tr.Sqlc.UpdateTopic(ctx, sqlc.UpdateTopicParams{
		ID:           int64(t.ID),
		RestaurantID: t.RestaurantID,
		Name:         t.Name,
		Slug:         &t.Slug,
		ParentID:     (*int64)(t.ParentID),
		SortOrder:    t.SortOrder,
	}); err != nil {
		return err
	}
	return nil
}

func (tr *TopicRepo) DeleteTopic(ctx context.Context, id domain.TopicID, restaurantID int32) error {
	return tr.Sqlc.DeleteTopic(ctx, sqlc.DeleteTopicParams{
		ID:           int64(id),
		RestaurantID: restaurantID,
	})
}
