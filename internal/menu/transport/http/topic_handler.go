package menuhttp

import (
	topicapp "go-ai/internal/menu/application/topic"
	"go-ai/internal/transport/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type TopicHandler struct {
	CreateUseCase *topicapp.CreateUseCase
	Logger        zerolog.Logger
}

func NewTopicHandler(createUseCase *topicapp.CreateUseCase, logger zerolog.Logger) *TopicHandler {
	return &TopicHandler{
		CreateUseCase: createUseCase,
		Logger:        logger,
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
	err := h.CreateUseCase.Execute(c.Request().Context(), in)
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
	// idInt, err := strconv.Atoi(id)
	// if err != nil {
	// 	return response.Error(c, http.StatusBadRequest, "invalid topic id format")
	// }
	// if idInt > math.MaxInt64 || idInt < math.MinInt64 {
	// 	return response.Error(c, http.StatusBadRequest, "topic id out of int64 range")
	// }
	return response.Success[any](c, nil, "Get topic successfully")
}
