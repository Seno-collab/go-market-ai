package menuhttp

import (
	menuitemapp "go-ai/internal/menu/application/menu_item"
	"go-ai/internal/transport/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const (
	errMissingTopicID       = "missing topic id"
	errInvalidTopicIDFormat = "invalid topic id format"
	errFailedGetMenuItem    = "Failed to get menu item"
	logGetMenuItemError     = "Get menu item error"
)

type MenuItemHandler struct {
	CreateUseCase       *menuitemapp.CreateUseCase
	GetUseCase          *menuitemapp.GetUseCase
	UpdateUseCase       *menuitemapp.UpdateUseCase
	DeleteUseCase       *menuitemapp.DeleteUseCase
	GetMenuItemsUseCase *menuitemapp.GetMenuItemsUseCase
	Logger              zerolog.Logger
}

func NewMenuItemHandler(
	createUseCase *menuitemapp.CreateUseCase,
	getUseCase *menuitemapp.GetUseCase,
	updateUseCase *menuitemapp.UpdateUseCase,
	deleteUseCase *menuitemapp.DeleteUseCase,
	getMenuItemsUseCase *menuitemapp.GetMenuItemsUseCase,
	logger zerolog.Logger) *MenuItemHandler {
	return &MenuItemHandler{
		CreateUseCase:       createUseCase,
		GetUseCase:          getUseCase,
		DeleteUseCase:       deleteUseCase,
		UpdateUseCase:       updateUseCase,
		GetMenuItemsUseCase: getMenuItemsUseCase,
		Logger:              logger,
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
// @Router /api/menu/items [post]
func (h *MenuItemHandler) Create(c echo.Context) error {
	var in menuitemapp.CreateMenuItemRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, errInvalidRequestPayload)
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	err = h.CreateUseCase.Execute(c.Request().Context(), in, restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Create menu item error")
		return response.Error(c, http.StatusBadRequest, "Failed to create menu item")
	}
	return response.Success[any](c, nil, "Create menu item successfully")
}

// GetMenuItemHandler godoc
// @Summary Get menu item by ID
// @Description Get a menu item detail by its ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Menu item ID"
// @Success 200 {object} app.GetMenuItemSuccessResponseDoc "Get menu item successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/items/{id} [get]
func (h *MenuItemHandler) Get(c echo.Context) error {
	idInt64, err := parseRequiredIDParam(c.Param("id"), errMissingTopicID, errInvalidTopicIDFormat)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	resp, err := h.GetUseCase.Execute(c.Request().Context(), idInt64, restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logGetMenuItemError)
		return response.Error(c, http.StatusBadRequest, errFailedGetMenuItem)
	}
	return response.Success(c, resp, "Get menu item successfully")
}

// UpdateMenuItemHandler godoc
// @Summary Update menu item
// @Description Update menu item by ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Menu item ID"
// @Param data body menuitemapp.UpdateMenuItemRequest true "Update menu item data"
// @Success 200 {object} app.UpdateMenuItemSuccessResponseDoc "Update menu item successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/items/{id} [put]
func (h *MenuItemHandler) Update(c echo.Context) error {
	idInt64, err := parseRequiredIDParam(c.Param("id"), errMissingTopicID, errInvalidTopicIDFormat)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	var in menuitemapp.UpdateMenuItemRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, errInvalidRequestPayload)
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	err = h.UpdateUseCase.Execute(c.Request().Context(), in, restaurantID, idInt64)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logGetMenuItemError)
		return response.Error(c, http.StatusBadRequest, errFailedGetMenuItem)
	}
	return response.Success[any](c, nil, "Update menu item successfully")
}

// DeleteMenuItemHandler godoc
// @Summary Delete menu item
// @Description Delete menu item by ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path string true "Menu item ID"
// @Success 200 {object} app.DeleteMenuItemSuccessResponseDoc "Delete menu item successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/items/{id} [delete]

func (h *MenuItemHandler) Delete(c echo.Context) error {
	idInt64, err := parseRequiredIDParam(c.Param("id"), errMissingTopicID, errInvalidTopicIDFormat)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	var in menuitemapp.UpdateMenuItemRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, errInvalidRequestPayload)
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	err = h.DeleteUseCase.Execute(c.Request().Context(), idInt64, restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logGetMenuItemError)
		return response.Error(c, http.StatusBadRequest, errFailedGetMenuItem)
	}
	return response.Success[any](c, nil, "Delete menu item successfully")
}

// GetMenuItemsByRestaurantHandler godoc
// @Summary Get menu items by restaurant
// @Description Get list of menu items by restaurant ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} app.GetMenuItemsByRestaurantSuccessResponseDoc "Get menu items by restaurant successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/restaurant/items [get]
func (h *MenuItemHandler) GetItems(c echo.Context) error {
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	resp, err := h.GetMenuItemsUseCase.Execute(c.Request().Context(), restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logGetMenuItemError)
		return response.Error(c, http.StatusBadRequest, errFailedGetMenuItem)
	}
	return response.Success(c, resp, "Get menu items by restaurant successfully")
}
