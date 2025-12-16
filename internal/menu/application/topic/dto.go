package topicapp

type CreateTopicRequest struct {
	RestaurantID int32  `json:"restaurant_id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ParentID     int64  `json:"parent_id"`
	SortOrder    int32  `json:"sort_order"`
}

type GetTopicRequest struct {
	ID string `json:"id"`
}

type GetTopicResponse struct {
	RestaurantID int32  `json:"restaurant_id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ParentID     int64  `json:"parent_id"`
	SortOrder    int32  `json:"sort_order"`
}
