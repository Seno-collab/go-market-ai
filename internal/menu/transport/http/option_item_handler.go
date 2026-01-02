package menuhttp

import (
	optionitemapp "go-ai/internal/menu/application/option_item"
	"go-ai/internal/transport/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type OptionItemHandler struct {
	CreateUseCase     *optionitemapp.CreateUseCase
	GetUseCase        *optionitemapp.GetUseCase
	GetByGroupUseCase *optionitemapp.GetByGroupUseCase
	UpdateUseCase     *optionitemapp.UpdateUseCase
	DeleteUseCase     *optionitemapp.DeleteUseCase
	Logger            zerolog.Logger
}

func NewOptionItemHandler(
	createUseCase *optionitemapp.CreateUseCase,
	getUseCase *optionitemapp.GetUseCase,
	getByGroupUseCase *optionitemapp.GetByGroupUseCase,
	updateUseCase *optionitemapp.UpdateUseCase,
	deleteUseCase *optionitemapp.DeleteUseCase,
	logger zerolog.Logger,
) *OptionItemHandler {
	return &OptionItemHandler{
		CreateUseCase:     createUseCase,
		GetUseCase:        getUseCase,
		GetByGroupUseCase: getByGroupUseCase,
		UpdateUseCase:     updateUseCase,
		DeleteUseCase:     deleteUseCase,
		Logger:            logger,
	}
}

// CreateOptionItem godoc
// @Summary Create option item
// @Description Create a new option item under an option group
// @Tags OptionItem
// @Accept json
// @Produce json
// @Param data body optionitemapp.CreateOptionItemRequest true "Option item payload"
// @Success 200 {object} app.CreateOptionItemSuccessResponseDoc "Create option item successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/option-item [post]
func (h *OptionItemHandler) Create(c echo.Context) error {
	var req optionitemapp.CreateOptionItemRequest
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
		h.Logger.Error().Err(err).Msg("create option item error")
		return response.Error(c, http.StatusBadRequest, "Failed to create option item")
	}
	return response.Success(c, &optionitemapp.CreateOptionItemResponse{ID: id}, "Create option item successfully")
}

// GetOptionItem godoc
// @Summary Get option item
// @Description Get option item details by ID
// @Tags OptionItem
// @Accept json
// @Produce json
// @Param id path string true "Option item ID"
// @Success 200 {object} app.GetOptionItemSuccessResponseDoc "Get option item successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/option-item/{id} [get]
func (h *OptionItemHandler) Get(c echo.Context) error {
	id, err := parseRequiredIDParam(c.Param("id"), "Invalid option item ID", "Invalid option item ID")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	item, err := h.GetUseCase.Execute(c.Request().Context(), id, restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("get option item error")
		return response.Error(c, http.StatusBadRequest, "Failed to get option item")
	}
	return response.Success(c, item, "Get option item successfully")
}

// GetOptionItemsByGroup godoc
// @Summary List option items in a group
// @Description Get option items for a specific option group
// @Tags OptionItem
// @Accept json
// @Produce json
// @Param id path string true "Option group ID"
// @Success 200 {object} app.GetOptionItemsSuccessResponseDoc "Get option items successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/option-group/{id}/option-items [get]
func (h *OptionItemHandler) GetByGroup(c echo.Context) error {
	var in optionitemapp.GetOptionItemsRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, errInvalidRequestPayload)
	}
	groupID, err := parseRequiredIDParam(c.Param("id"), "Invalid option group ID", "Invalid option group ID")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	resp, err := h.GetByGroupUseCase.Execute(c.Request().Context(), groupID, restaurantID, in)
	if err != nil {
		h.Logger.Error().Err(err).Msg("get option items error")
		return response.Error(c, http.StatusBadRequest, "Failed to get option items")
	}
	return response.Success(c, resp, "Get option items successfully")
}

// UpdateOptionItem godoc
// @Summary Update option item
// @Description Update option item details by ID
// @Tags OptionItem
// @Accept json
// @Produce json
// @Param id path string true "Option item ID"
// @Param data body optionitemapp.UpdateOptionItemRequest true "Option item payload"
// @Success 200 {object} app.UpdateOptionItemSuccessResponseDoc "Update option item successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/option-item/{id} [put]
func (h *OptionItemHandler) Update(c echo.Context) error {
	id, err := parseRequiredIDParam(c.Param("id"), "Invalid option item ID", "Invalid option item ID")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	var req optionitemapp.UpdateOptionItemRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errInvalidRequestPayload)
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	if err := h.UpdateUseCase.Execute(c.Request().Context(), id, req, restaurantID); err != nil {
		h.Logger.Error().Err(err).Msg("update option item error")
		return response.Error(c, http.StatusBadRequest, "Failed to update option item")
	}
	return response.Success[any](c, nil, "Update option item successfully")
}

// DeleteOptionItem godoc
// @Summary Delete option item
// @Description Delete option item by ID
// @Tags OptionItem
// @Accept json
// @Produce json
// @Param id path string true "Option item ID"
// @Success 200 {object} app.DeleteOptionItemSuccessResponseDoc "Delete option item successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/option-item/{id} [delete]
func (h *OptionItemHandler) Delete(c echo.Context) error {
	id, err := parseRequiredIDParam(c.Param("id"), "Invalid option item ID", "Invalid option item ID")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	if err := h.DeleteUseCase.Execute(c.Request().Context(), id, restaurantID); err != nil {
		h.Logger.Error().Err(err).Msg("delete option item error")
		return response.Error(c, http.StatusBadRequest, "Failed to delete option item")
	}
	return response.Success[any](c, nil, "Delete option item successfully")
}
