package optiongroupapp

type CreateOptionGroupRequest struct {
	MenuItemID int64  `json:"menu_item_id"`
	Name       string `json:"name"`
	MinSelect  int32  `json:"min_select"`
	MaxSelect  int32  `json:"max_select"`
	IsRequired bool   `json:"is_required"`
	SortOrder  int32  `json:"sort_order"`
}

type CreateOptionGroupResponse struct {
	ID int64 `json:"id"`
}

type UpdateOptionGroupRequest struct {
	Name       string `json:"name"`
	MinSelect  int32  `json:"min_select"`
	MaxSelect  int32  `json:"max_select"`
	IsRequired bool   `json:"is_required"`
	SortOrder  int32  `json:"sort_order"`
}

type GetOptionGroupResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	MinSelect  int32  `json:"min_select"`
	MaxSelect  int32  `json:"max_select"`
	IsRequired bool   `json:"is_required"`
	SortOrder  int32  `json:"sort_order"`
}

type GetOptionGroupsResponse struct {
	Groups []GetOptionGroupResponse `json:"groups"`
}
