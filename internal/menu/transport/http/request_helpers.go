package menuhttp

import (
	"errors"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	errInvalidRequestPayload = "Invalid request payload"
	errInvalidRestaurantID   = "Invalid restaurantID"
	logInvalidRestaurantID   = "invalid restaurant id"
)

var (
	errRestaurantIDMissing = errors.New("missing restaurant id")
	errRestaurantIDInvalid = errors.New("invalid restaurant id")
)

func parseRequiredIDParam(raw, missingMsg, invalidMsg string) (int64, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, errors.New(missingMsg)
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, errors.New(invalidMsg)
	}
	return id, nil
}

func getRestaurantID(c echo.Context) (int32, error) {
	val := c.Get("restaurant_id")
	if val == nil {
		return 0, errRestaurantIDMissing
	}
	restaurantID, ok := val.(int32)
	if !ok {
		return 0, errRestaurantIDInvalid
	}
	return restaurantID, nil
}
