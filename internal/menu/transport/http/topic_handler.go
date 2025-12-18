package menuhttp

import (
	topicapp "go-ai/internal/menu/application/topic"
	"go-ai/internal/transport/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type TopicHandler struct {
	CreateUseCase    *topicapp.CreateUseCase
	GetUseCase       *topicapp.GetUseCase
	GetTopicsUseCase *topicapp.GetTopicsUseCase
	Logger           zerolog.Logger
}

func NewTopicHandler(createUseCase *topicapp.CreateUseCase, getUseCase *topicapp.GetUseCase, getTopicesUseCase *topicapp.GetTopicsUseCase, logger zerolog.Logger) *TopicHandler {
	return &TopicHandler{
		CreateUseCase:    createUseCase,
		GetUseCase:       getUseCase,
		GetTopicsUseCase: getTopicesUseCase,
		Logger:           logger,
	}
}

// CreateTopicHandler godoc
// @Summary      Create a new topic
// @Description  Create a new topic with name, price, description, and optional image or logo
// @Tags         Topic
// @Accept       json
// @Produce      json
// @Param        data  body      topicapp.CreateTopicRequest  true  "Topic data"
// @Success      200   {object}  app.CreateTopicSuccessResponseDoc  "Create topic success"
// @Failure      default {object} response.ErrorDoc                "Errors"
// @Router       /api/menu/topic [post]
func (h *TopicHandler) Create(c echo.Context) error {
	var in topicapp.CreateTopicRequest
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
		h.Logger.Error().Err(err).Msg("Create topic error")
		return response.Error(c, http.StatusBadRequest, "Failed to create topic")
	}
	return response.Success[any](c, nil, "Create topic successfully")
}

// GetRestaurant godoc
// @Summary Get restaurant by ID
// @Description Get detailed information of a restaurant using its ID
// @Tags Restaurant
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Success 200 {object} app.GetTopicSuccessResponseDoc "Get topic successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/topic/{id} [get]
func (h *TopicHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "missing topic id")
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid topic id format")
	}
	restaurantID := c.Get("restaurant_id")
	if restaurantID == nil {
		h.Logger.Error().Msg("restaurantID is nil")
		return response.Error(c, http.StatusBadRequest, "Invalid restaurantID")
	}
	resp, err := h.GetUseCase.Execute(c.Request().Context(), idInt64, restaurantID.(int32))
	if err != nil {
		h.Logger.Error().Err(err).Msg("Get topic error")
		return response.Error(c, http.StatusBadRequest, "Failed to get topic")
	}
	return response.Success(c, resp, "Get topic successfully")
}

// GetTopicsByRestaurantHandler godoc
// @Summary Get topics by restaurant
// @Description Get list of topics by restaurant ID
// @Tags Topic
// @Accept json
// @Produce json
// @Param restaurant_id path string true "Restaurant ID"
// @Success 200 {object} app.GetTopicsByRestaurantSuccessResponseDoc "Get topics by restaurant successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/restaurant/topics [get]
func (h *TopicHandler) GetTopcis(c echo.Context) error {
	restaurantID := c.Get("restaurant_id")
	if restaurantID == nil {
		h.Logger.Error().Msg("restaurantID is nil")
		return response.Error(c, http.StatusBadRequest, "Invalid restaurantID")
	}
	resp, err := h.GetTopicsUseCase.Execute(c.Request().Context(), restaurantID.(int32))
	if err != nil {
		h.Logger.Error().Err(err).Msg("Get topic error")
		return response.Error(c, http.StatusBadRequest, "Failed to get topics")
	}
	return response.Success(c, resp, "Get topics successfully")
}
