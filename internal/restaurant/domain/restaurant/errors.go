package restaurant

import (
	domainerr "go-ai/pkg/domain_err"
	"net/http"
)

var (
	ErrNameRequired        = domainerr.New(http.StatusBadRequest, "name is required")
	ErrDescriptionRequired = domainerr.New(http.StatusBadRequest, "description is required")
	ErrAddressRequired     = domainerr.New(http.StatusBadRequest, "address is required")
	ErrCategoryRequired    = domainerr.New(http.StatusBadRequest, "category is required")
	ErrCityRequired        = domainerr.New(http.StatusBadRequest, "city is required")
	ErrDistrictRequired    = domainerr.New(http.StatusBadRequest, "district is required")

	ErrInvalidPhone   = domainerr.New(http.StatusBadRequest, "invalid phone number")
	ErrInvalidWebsite = domainerr.New(http.StatusBadRequest, "invalid website url")
	ErrUserIDRequired = domainerr.New(http.StatusBadRequest, "userID is required")

	ErrHoursInvalidDay     = domainerr.New(http.StatusBadRequest, "invalid day of week")
	ErrHoursInvalidTime    = domainerr.New(http.StatusBadRequest, "open time must be before close time")
	ErrHoursTimeFormat     = domainerr.New(http.StatusBadRequest, "invalid time format (must be HH:mm)")
	ErrInvalidUrl          = domainerr.New(http.StatusBadRequest, "invalid url")
	ErrRestaurantNotExists = domainerr.New(http.StatusNotFound, "restaurant does not exist")
)
