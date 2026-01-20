package menuhttp

import (
	"net/http"
	"strconv"
	"strings"

	menuapp "go-ai/internal/menu/application/menu"
	"go-ai/internal/menu/domain"
	"go-ai/internal/transport/response"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const (
	errInvalidMenuType = "Invalid type. Allowed values: dish, beverage, extra, combo"
	errLimitMustBePos  = "limit must be a positive integer"
	errCursorInvalid   = "cursor must be a valid int64"
)

type MenuHandler struct {
	ListUseCase *menuapp.ListMenusUseCase
	Logger      zerolog.Logger
}

func NewMenuHandler(listUseCase *menuapp.ListMenusUseCase, logger zerolog.Logger) *MenuHandler {
	return &MenuHandler{
		ListUseCase: listUseCase,
		Logger:      logger.With().Str("handler", "MenuHandler").Logger(),
	}
}

// ListMenus godoc
// @Summary      List menus with cursor pagination
// @Description  Retrieve menus filtered by restaurant, type, and optional topics using cursor-based pagination
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        X-Restaurant-Id header int true "Restaurant ID"
// @Param        type query string true "Menu type" Enums(dish,beverage,extra,combo)
// @Param        topics query string false "Comma-separated topic names, e.g. \"2024,2025\""
// @Param        limit query int false "Items per page (default 20, max 50)"
// @Param        cursor query int64 false "Pagination cursor (last returned menu id)"
// @Success      200 {object} menuapp.ListMenusResponse "Menus page"
// @Failure      default {object} response.ErrorDoc "Errors"
// @Router       /api/menus [get]
func (h *MenuHandler) ListMenus(c echo.Context) error {
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}

	menuType := domain.MenuType(strings.TrimSpace(c.QueryParam("type")))
	if !menuType.Valid() && menuType != "" {
		return response.Error(c, http.StatusBadRequest, errInvalidMenuType)
	}

	rawTopics := strings.TrimSpace(c.QueryParam("topics"))
	topics := make([]string, 0)
	if rawTopics != "" {
		seen := make(map[string]struct{})
		for remaining := rawTopics; ; {
			segment, rest, found := strings.Cut(remaining, ",")
			if trimmed := strings.TrimSpace(segment); trimmed != "" {
				if _, ok := seen[trimmed]; !ok {
					seen[trimmed] = struct{}{}
					topics = append(topics, trimmed)
				}
			}
			if !found {
				break
			}
			remaining = rest
		}
	}

	var limit int32
	if rawLimit := strings.TrimSpace(c.QueryParam("limit")); rawLimit != "" {
		l, err := strconv.Atoi(rawLimit)
		if err != nil || l < 1 {
			return response.Error(c, http.StatusBadRequest, errLimitMustBePos)
		}
		limit = int32(l)
	}

	var cursor *int64
	if rawCursor := strings.TrimSpace(c.QueryParam("cursor")); rawCursor != "" {
		val, err := strconv.ParseInt(rawCursor, 10, 64)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, errCursorInvalid)
		}
		cursor = &val
	}

	resp, err := h.ListUseCase.Execute(c.Request().Context(), menuapp.ListMenusRequest{
		RestaurantID: restaurantID,
		Type:         menuType,
		Topics:       topics,
		Limit:        limit,
		Cursor:       cursor,
	})
	if err != nil {
		h.Logger.Error().Err(err).Msg("List menus error")
		return response.Error(c, http.StatusInternalServerError, "Failed to get menus")
	}
	return response.Success(c, resp, "Get menus successfully")
}
