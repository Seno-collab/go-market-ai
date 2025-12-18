package domain

import "errors"

var (
	ErrNameRequired             = errors.New("menu: name is required")
	ErrInvalidPrice             = errors.New("menu: price must be >= 0")
	ErrInvalidVariantDelta      = errors.New("variant: price delta must be >= 0")
	ErrVariantDefaultMultiple   = errors.New("variant: only one default variant allowed")
	ErrOptionGroupInvalidMinMax = errors.New("option_group: max_select must be >= min_select")
	ErrOptionItemInvalidPrice   = errors.New("option_item: price delta must >= 0")
	ErrComboNotComboType        = errors.New("combo: must be type=combo")
	ErrRecordUpdateFailed       = errors.New("record update failed")
)
