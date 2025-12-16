package menuhttp

import (
	menuitemapp "go-ai/internal/menu/application/menu_item"
	"go-ai/internal/transport/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type MenuItemHandler struct {
	CreateUseCase *menuitemapp.CreateUseCase
	Logger        zerolog.Logger
}

func NewMenuItemHandler(createUseCase *menuitemapp.CreateUseCase, logger zerolog.Logger) *MenuItemHandler {
	return &MenuItemHandler{
		CreateUseCase: createUseCase,
		Logger:        logger,
	}
}

// CreateMenuItemHandler godoc
// @Summary Create a new menu item
// @Description Create a new menu item with name, price, description and optional image/logo
// @Tags Menu
// @Accept json
// @Produce json
// @Param data body menuitemapp.CreateMenuItemRequest true "Menu item data"
// @Success 200 {object} app.CreateMenuItemSuccessResponseDoc "Create menu item success"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/item [post]
func (h *MenuItemHandler) Create(c echo.Context) error {
	var in menuitemapp.CreateMenuItemRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	restaurantID := c.Get("restaurant_id")
	if restaurantID == nil {
		h.Logger.Error().Msg("restaurantID is nil")
		return response.Error(c, http.StatusBadRequest, "Invalid restaurantID")
	}
	err := h.CreateUseCase.Execute(c.Request().Context(), in, restaurantID.(int32))
	if err != nil {
		h.Logger.Error().Err(err).Msg("Create menu item error")
		return response.Error(c, http.StatusBadRequest, "Failed to create menu item")
	}
	return response.Success[any](c, nil, "Create menu item successfully")
}
