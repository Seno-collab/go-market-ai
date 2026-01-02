package domain

import "errors"

var (
	ErrNameRequired              = errors.New("Menu: name is required")
	ErrInvalidPrice              = errors.New("Menu: price must be >= 0")
	ErrInvalidVariantDelta       = errors.New("Variant: price delta must be >= 0")
	ErrVariantDefaultMultiple    = errors.New("Variant: only one default variant allowed")
	ErrOptionGroupEmptyName      = errors.New("Option_group: name cannot be empty")
	ErrOptionGroupNegativeValues = errors.New("Option_group: min_select and max_select must be >= 0")
	ErrOptionGroupInvalidMinMax  = errors.New("Option_group: max_select must be >= min_select")
	ErrOptionItemInvalidPrice    = errors.New("Option_item: price delta must >= 0")
	ErrComboNotComboType         = errors.New("Combo: must be type=combo")
	ErrRecordUpdateFailed        = errors.New("Record update failed")
)
