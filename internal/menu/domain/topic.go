package domain

import "errors"

var (
	ErrTopicNameRequired = errors.New("Topic: name is required")
	ErrTopicInvalidLevel = errors.New("Topic: circular parent is not allowed")
)

type TopicID int64

type Topic struct {
	ID           TopicID
	RestaurantID int32
	Name         string
	Slug         string
	ParentID     *TopicID
	SortOrder    int32
	IsActive     bool
}

// Factory
func NewTopic(restaurantID int32, name string, slug string, parent *TopicID, sort int32) (*Topic, error) {
	if name == "" {
		return nil, ErrTopicNameRequired
	}

	return &Topic{
		RestaurantID: restaurantID,
		Name:         name,
		Slug:         slug,
		ParentID:     parent,
		SortOrder:    sort,
		IsActive:     true,
	}, nil
}

// Validation
func (t *Topic) Validate() error {
	if t.Name == "" {
		return ErrTopicNameRequired
	}
	if t.SortOrder < 0 {
		return errors.New("Topic: invalid sort order")
	}
	return nil
}

// Business rule: prevent self-parent
func (t *Topic) ValidateParent() error {
	if t.ParentID != nil && *t.ParentID == t.ID {
		return ErrTopicInvalidLevel
	}
	return nil
}

// Activate/Deactivate
func (t *Topic) Deactivate() {
	t.IsActive = false
}

func (t *Topic) Activate() {
	t.IsActive = true
}
