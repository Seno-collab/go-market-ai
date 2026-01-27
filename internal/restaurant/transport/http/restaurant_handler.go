package restauranthttp

import (
	restaurantapp "go-ai/internal/restaurant/application/restaurant"
	"go-ai/internal/transport/response"
	domainerr "go-ai/pkg/domain_err"
	"math"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
)

type RestaurantHandler struct {
	CreateUseCase      *restaurantapp.CreateRestaurantUseCase
	GetByIDUseCase     *restaurantapp.GetByIDUseCase
	UpdateUseCase      *restaurantapp.UpdateRestaurantUseCase
	DeleteUseCase      *restaurantapp.DeleteUseCase
	GetComboboxUseCase *restaurantapp.GetRestaurantItemComboboxUseCase
	Logger             zerolog.Logger
}

func NewRestaurantHandler(
	createUseCase *restaurantapp.CreateRestaurantUseCase,
	getByIDUseCase *restaurantapp.GetByIDUseCase,
	updateUseCase *restaurantapp.UpdateRestaurantUseCase,
	deleteUseCase *restaurantapp.DeleteUseCase,
	getComboboxUseCase *restaurantapp.GetRestaurantItemComboboxUseCase,
	logger zerolog.Logger,
) *RestaurantHandler {
	return &RestaurantHandler{
		CreateUseCase:      createUseCase,
		GetByIDUseCase:     getByIDUseCase,
		UpdateUseCase:      updateUseCase,
		DeleteUseCase:      deleteUseCase,
		GetComboboxUseCase: getComboboxUseCase,
		Logger:             logger.With().Str("component", "RestaurantHandler").Logger(),
	}
}

// CreateRestaurant godoc
// @Summary Create restaurant
// @Description Create a new restaurant with name, email, phone, logo_url, banner_url,...
// @Tags Restaurant
// @Accept json
// @Produce json
// @Param request body restaurantapp.CreateRestaurantRequest true "Restaurant create payload"
// @Success 200 {object} app.CreateRestaurantSuccessResponseDoc "Create restaurant successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/restaurants [post]
func (h *RestaurantHandler) Create(c *echo.Context) error {
	var in restaurantapp.CreateRestaurantRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	userID := c.Get("user_id")
	if userID == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	id, err := h.CreateUseCase.Execute(c.Request().Context(), in, userUUID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed create restaurant")
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success(c, restaurantapp.CreateRestaurantResponse{
		ID: id,
	}, "Create restaurant successfully")
}

// GetRestaurant godoc
// @Summary Get restaurant by ID
// @Description Get detailed information of a restaurant using its ID
// @Tags Restaurant
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} app.GetRestaurantByIDSuccessResponseDoc "Get restaurant successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/restaurants/{id} [get]
func (h *RestaurantHandler) GetByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "Missing restaurant id")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid restaurant id format")
	}
	if idInt > math.MaxInt32 || idInt < math.MinInt32 {
		return response.Error(c, http.StatusBadRequest, "Restaurant id out of int32 range")
	}
	restaurant, err := h.GetByIDUseCase.Execute(c.Request().Context(), int32(idInt))
	if restaurant == nil {
		return response.Error(c, http.StatusNotFound, "Restaurant not found")
	}
	return response.Success(c, restaurant, "Get restaurant successfully")
}

// UpdateRestaurant godoc
// @Summary Update restaurant information
// @Description Update restaurant fields such as name, address, contact info, logo, banner, etc.
// @Tags Restaurant
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Param body body restaurantapp.UpdateRestaurantRequest true "Restaurant update payload"
// @Success 200 {object} app.UpdateRestaurantSuccessResponseDoc "Update restaurant successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/restaurants/{id} [put]
func (h *RestaurantHandler) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "Missing restaurant id")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid restaurant id format")
	}
	if idInt > math.MaxInt32 || idInt < math.MinInt32 {
		return response.Error(c, http.StatusBadRequest, "Restaurant id out of int32 range")
	}
	var in restaurantapp.CreateRestaurantRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	userId := c.Get("user_id")
	if userId == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userId.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}

	if err := h.UpdateUseCase.Execute(c.Request().Context(), in, userUUID, int32(idInt)); err != nil {
		h.Logger.Error().Err(err).Msg("failed create restaurant")
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success[any](c, nil, "Create restaurant successfully")
}

// DeleteRestaurant godoc
// @Summary Delete restaurant by ID
// @Description Remove a restaurant and its related data using its ID
// @Tags Restaurant
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} app.DeleteRestaurantSuccessResponseDoc "Restaurant deleted successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/restaurants/{id} [delete]
func (h *RestaurantHandler) Delete(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "Missing restaurant id")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid restaurant id format")
	}
	if idInt > math.MaxInt32 || idInt < math.MinInt32 {
		return response.Error(c, http.StatusBadRequest, "Restaurant id out of int32 range")
	}
	userID := c.Get("user_id")
	if userID == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	err = h.DeleteUseCase.Execute(c.Request().Context(), int32(idInt), userUUID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Delete restaurant failed")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success[any](c, nil, "Restaurant deleted successfully")
}

// GetRestaurantCombobox godoc
// @Summary      Get combobox of a menu item
// @Description  Get option/combobox groups and items
//
//	using the restaurant context selected by the current authenticated user
//
// @Tags         Restaurant
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} app.GetComboboxRestaurantSuccessResponseDoc "Get restaurant combobox successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router       /api/restaurants/combobox [get]
func (h *RestaurantHandler) GetCombobox(c *echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	resp, err := h.GetComboboxUseCase.Execute(c.Request().Context(), userUUID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Get restaurant combobox failed")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success(c, resp, "Get restaurant combobox successfully")
}
