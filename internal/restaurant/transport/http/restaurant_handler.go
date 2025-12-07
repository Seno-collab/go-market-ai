package restauranthttp

import (
	restaurantapp "go-ai/internal/restaurant/application/restaurant"
	"go-ai/internal/transport/response"
	domainerr "go-ai/pkg/domain_err"
	"math"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type RestaurantHandler struct {
	CreateUC  *restaurantapp.CreateRestaurantUseCase
	GetByIdUC *restaurantapp.GetByIDUseCase
	UpdateUC  *restaurantapp.UpdateRestaurantUseCase
	DeleteUC  *restaurantapp.DeleteUseCase
	Logger    zerolog.Logger
}

func NewRestaurantHandler(
	createUC *restaurantapp.CreateRestaurantUseCase,
	getByIDUC *restaurantapp.GetByIDUseCase,
	updateUC *restaurantapp.UpdateRestaurantUseCase,
	deleteUC *restaurantapp.DeleteUseCase,
	logger zerolog.Logger) *RestaurantHandler {
	return &RestaurantHandler{
		CreateUC:  createUC,
		GetByIdUC: getByIDUC,
		UpdateUC:  updateUC,
		DeleteUC:  deleteUC,
		Logger:    logger.With().Str("component", "Restaurant handler").Logger(),
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
// @Router /api/restaurant [post]
func (h *RestaurantHandler) Create(c echo.Context) error {
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
	id, err := h.CreateUC.Execute(c.Request().Context(), in, userUUID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed create restaurant")
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success[restaurantapp.CreateRestaurantResponse](c, &restaurantapp.CreateRestaurantResponse{
		Id: id,
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
// @Router /api/restaurant/{id} [get]
func (h *RestaurantHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "missing restaurant id")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid restaurant id format")
	}
	if idInt > math.MaxInt32 || idInt < math.MinInt32 {
		return response.Error(c, http.StatusBadRequest, "restaurant id out of int32 range")
	}
	restaurant, err := h.GetByIdUC.Execute(c.Request().Context(), int32(idInt))
	if restaurant == nil {
		return response.Error(c, http.StatusNotFound, "restaurant not found")
	}
	return response.Success[restaurantapp.GetRestaurantByIDResponse](c, restaurant, "Get restaurant successfully")
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
// @Router /api/restaurant/{id} [put]
func (h *RestaurantHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "missing restaurant id")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid restaurant id format")
	}
	if idInt > math.MaxInt32 || idInt < math.MinInt32 {
		return response.Error(c, http.StatusBadRequest, "restaurant id out of int32 range")
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

	if err := h.UpdateUC.Execute(c.Request().Context(), in, userUUID, int32(idInt)); err != nil {
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
// @Router /api/restaurant/{id} [delete]
func (h *RestaurantHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "missing restaurant id")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid restaurant id format")
	}
	if idInt > math.MaxInt32 || idInt < math.MinInt32 {
		return response.Error(c, http.StatusBadRequest, "restaurant id out of int32 range")
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
	err = h.DeleteUC.Execute(c.Request().Context(), int32(idInt), userUUID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Delete restaurant failed")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success[any](c, nil, "Restaurant deleted successfully")
}
