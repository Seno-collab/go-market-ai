package restaurant

import (
	domainerr "go-ai/pkg/domain_err"
	"net/http"
)

var (
	ErrNameRequired        = domainerr.New(http.StatusBadRequest, "Name is required")
	ErrDescriptionRequired = domainerr.New(http.StatusBadRequest, "Description is required")
	ErrAddressRequired     = domainerr.New(http.StatusBadRequest, "Address is required")
	ErrCategoryRequired    = domainerr.New(http.StatusBadRequest, "Category is required")
	ErrCityRequired        = domainerr.New(http.StatusBadRequest, "City is required")
	ErrDistrictRequired    = domainerr.New(http.StatusBadRequest, "District is required")

	ErrInvalidPhone   = domainerr.New(http.StatusBadRequest, "Invalid phone number")
	ErrInvalidWebsite = domainerr.New(http.StatusBadRequest, "Invalid website url")
	ErrUserIDRequired = domainerr.New(http.StatusBadRequest, "UserID is required")

	ErrHoursInvalidDay     = domainerr.New(http.StatusBadRequest, "Invalid day of week")
	ErrHoursInvalidTime    = domainerr.New(http.StatusBadRequest, "Open time must be before close time")
	ErrHoursTimeFormat     = domainerr.New(http.StatusBadRequest, "Invalid time format (must be HH:mm)")
	ErrRestaurantNotExists = domainerr.New(http.StatusNotFound, "Restaurant does not exist")
	ErrRestaurantExists    = domainerr.New(http.StatusBadRequest, "Restaurant exist")
)
