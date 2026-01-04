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
	val := c.Request().Header.Get("X-Restaurant-Id")
	if val == "" {
		if ctxVal := c.Get("restaurant_id"); ctxVal != nil {
			switch v := ctxVal.(type) {
			case int32:
				return v, nil
			case int:
				return int32(v), nil
			case int64:
				return int32(v), nil
			case string:
				restaurantID, err := strconv.ParseInt(v, 10, 32)
				if err == nil {
					return int32(restaurantID), nil
				}
			}
		}
		return 0, errRestaurantIDMissing
	}
	restaurantID, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, errRestaurantIDInvalid
	}
	return int32(restaurantID), nil
}
