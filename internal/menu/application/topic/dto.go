package topicapp

import "go-ai/internal/transport/response"

type CreateTopicRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ParentID  int64  `json:"parent_id"`
	SortOrder int32  `json:"sort_order"`
}

type UpdateTopicRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ParentID  int64  `json:"parent_id"`
	SortOrder int32  `json:"sort_order"`
}

type GetTopicRequest struct {
	ID string `json:"id"`
}

type GetTopicResponse struct {
	ID           int64  `json:"id"`
	RestaurantID int32  `json:"restaurant_id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ParentID     *int64 `json:"parent_id,omitempty"`
	SortOrder    int32  `json:"sort_order"`
}

type GetTopicsResponse struct {
	response.PaginatedResponse[[]GetTopicResponse]
}
