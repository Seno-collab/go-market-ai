package menuhttp

import (
	optiongroupapp "go-ai/internal/menu/application/option-group"
	"go-ai/internal/transport/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const errInvalidOptionGroupID = "Invalid option group ID"

type OptionGroupHandler struct {
	CreateUseCase        *optiongroupapp.CreateUseCase
	GetUseCase           *optiongroupapp.GetUseCase
	GetByMenuItemUseCase *optiongroupapp.GetByMenuItemUseCase
	UpdateUseCase        *optiongroupapp.UpdateUseCase
	DeleteUseCase        *optiongroupapp.DeleteUseCase
	Logger               zerolog.Logger
}

func NewOptionGroupHandler(
	createUseCase *optiongroupapp.CreateUseCase,
	getUseCase *optiongroupapp.GetUseCase,
	getByMenuItemUseCase *optiongroupapp.GetByMenuItemUseCase,
	updateUseCase *optiongroupapp.UpdateUseCase,
	deleteUseCase *optiongroupapp.DeleteUseCase,
	logger zerolog.Logger,
) *OptionGroupHandler {
	return &OptionGroupHandler{
		CreateUseCase:        createUseCase,
		GetUseCase:           getUseCase,
		GetByMenuItemUseCase: getByMenuItemUseCase,
		UpdateUseCase:        updateUseCase,
		DeleteUseCase:        deleteUseCase,
		Logger:               logger,
	}
}

// CreateOptionGroup godoc
// @Summary Create option group
// @Description Create a new option group for a menu item
// @Tags OptionGroup
// @Accept json
// @Produce json
// @Param data body optiongroupapp.CreateOptionGroupRequest true "Option group payload"
// @Success 200 {object} app.CreateOptionGroupSuccessResponseDoc "Create option group successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/option-group [post]
func (h *OptionGroupHandler) Create(c echo.Context) error {
	var req optiongroupapp.CreateOptionGroupRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errInvalidRequestPayload)
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	id, err := h.CreateUseCase.Execute(c.Request().Context(), req, restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("create option group error")
		return response.Error(c, http.StatusBadRequest, "Failed to create option group")
	}
	return response.Success(c, &optiongroupapp.CreateOptionGroupResponse{ID: id}, "Create option group successfully")
}

// GetOptionGroup godoc
// @Summary Get option group by ID
// @Description Get detailed information of an option group
// @Tags OptionGroup
// @Accept json
// @Produce json
// @Param id path string true "Option group ID"
// @Success 200 {object} app.GetOptionGroupSuccessResponseDoc "Get option group successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/option-group/{id} [get]
func (h *OptionGroupHandler) Get(c echo.Context) error {
	id, err := parseRequiredIDParam(c.Param("id"), errInvalidOptionGroupID, errInvalidOptionGroupID)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	group, err := h.GetUseCase.Execute(c.Request().Context(), id, restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("get option group error")
		return response.Error(c, http.StatusBadRequest, "Failed to get option group")
	}
	return response.Success(c, group, "Get option group successfully")
}

// GetOptionGroupsByMenuItem godoc
// @Summary List option groups of a menu item
// @Description Get option groups attached to a specific menu item
// @Tags OptionGroup
// @Accept json
// @Produce json
// @Param id path string true "Menu item ID"
// @Success 200 {object} app.GetOptionGroupsSuccessResponseDoc "Get option groups successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/item/{id}/option-groups [get]
func (h *OptionGroupHandler) GetByMenuItem(c echo.Context) error {
	menuItemID, err := parseRequiredIDParam(c.Param("id"), "Invalid menu item ID", "Invalid menu item ID")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	resp, err := h.GetByMenuItemUseCase.Execute(c.Request().Context(), menuItemID, restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("get option groups error")
		return response.Error(c, http.StatusBadRequest, "Failed to get option groups")
	}
	return response.Success(c, resp, "Get option groups successfully")
}

// UpdateOptionGroup godoc
// @Summary Update option group
// @Description Update option group details by ID
// @Tags OptionGroup
// @Accept json
// @Produce json
// @Param id path string true "Option group ID"
// @Param data body optiongroupapp.UpdateOptionGroupRequest true "Option group payload"
// @Success 200 {object} app.UpdateOptionGroupSuccessResponseDoc "Update option group successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/option-group/{id} [put]
func (h *OptionGroupHandler) Update(c echo.Context) error {
	id, err := parseRequiredIDParam(c.Param("id"), errInvalidOptionGroupID, errInvalidOptionGroupID)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	var req optiongroupapp.UpdateOptionGroupRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errInvalidRequestPayload)
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	if err := h.UpdateUseCase.Execute(c.Request().Context(), id, req, restaurantID); err != nil {
		h.Logger.Error().Err(err).Msg("update option group error")
		return response.Error(c, http.StatusBadRequest, "Failed to update option group")
	}
	return response.Success[any](c, nil, "Update option group successfully")
}

// DeleteOptionGroup godoc
// @Summary Delete option group
// @Description Delete an option group by ID
// @Tags OptionGroup
// @Accept json
// @Produce json
// @Param id path string true "Option group ID"
// @Success 200 {object} app.DeleteOptionGroupSuccessResponseDoc "Delete option group successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/option-group/{id} [delete]
func (h *OptionGroupHandler) Delete(c echo.Context) error {
	id, err := parseRequiredIDParam(c.Param("id"), errInvalidOptionGroupID, errInvalidOptionGroupID)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	if err := h.DeleteUseCase.Execute(c.Request().Context(), id, restaurantID); err != nil {
		h.Logger.Error().Err(err).Msg("delete option group error")
		return response.Error(c, http.StatusBadRequest, "Failed to delete option group")
	}
	return response.Success[any](c, nil, "Delete option group successfully")
}
