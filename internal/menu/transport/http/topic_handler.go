package menuhttp

import (
	topicapp "go-ai/internal/menu/application/topic"
	"go-ai/internal/transport/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type TopicHandler struct {
	CreateUseCase    *topicapp.CreateUseCase
	GetUseCase       *topicapp.GetUseCase
	GetTopicsUseCase *topicapp.GetTopicsUseCase
	UpdateUseCase    *topicapp.UpdateUseCase
	DeleteUseCase    *topicapp.DeleteUseCase
	Logger           zerolog.Logger
}

func NewTopicHandler(
	createUseCase *topicapp.CreateUseCase,
	getUseCase *topicapp.GetUseCase,
	getTopicsUseCase *topicapp.GetTopicsUseCase,
	updateUseCase *topicapp.UpdateUseCase,
	deleteUseCase *topicapp.DeleteUseCase,
	logger zerolog.Logger,
) *TopicHandler {
	return &TopicHandler{
		CreateUseCase:    createUseCase,
		GetUseCase:       getUseCase,
		GetTopicsUseCase: getTopicsUseCase,
		UpdateUseCase:    updateUseCase,
		DeleteUseCase:    deleteUseCase,
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
		return response.Error(c, http.StatusBadRequest, errInvalidRequestPayload)
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	err = h.CreateUseCase.Execute(c.Request().Context(), in, restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Create topic error")
		return response.Error(c, http.StatusBadRequest, "Failed to create topic")
	}
	return response.Success[any](c, nil, "Create topic successfully")
}

// GetTopic godoc
// @Summary Get topic by ID
// @Description Get detailed information of a topic using its ID
// @Tags Topic
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Success 200 {object} app.GetTopicSuccessResponseDoc "Get topic successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/topic/{id} [get]
func (h *TopicHandler) Get(c echo.Context) error {
	idInt64, err := parseRequiredIDParam(c.Param("id"), "missing topic id", "invalid topic id format")
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
		h.Logger.Error().Err(err).Msg("Get topic error")
		return response.Error(c, http.StatusBadRequest, "Failed to get topic")
	}
	return response.Success(c, resp, "Get topic successfully")
}

// UpdateTopicHandler godoc
// @Summary Update topic
// @Description Update topic information by ID
// @Tags Topic
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Param data body topicapp.UpdateTopicRequest true "Topic data"
// @Success 200 {object} app.UpdateTopicSuccessResponseDoc "Update topic successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/topic/{id} [put]
func (h *TopicHandler) Update(c echo.Context) error {
	idInt64, err := parseRequiredIDParam(c.Param("id"), "missing topic id", "invalid topic id format")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	var in topicapp.UpdateTopicRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, errInvalidRequestPayload)
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	if err := h.UpdateUseCase.Execute(c.Request().Context(), idInt64, in, restaurantID); err != nil {
		h.Logger.Error().Err(err).Msg("Update topic error")
		return response.Error(c, http.StatusBadRequest, "Failed to update topic")
	}
	return response.Success[any](c, nil, "Update topic successfully")
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
func (h *TopicHandler) GetTopics(c echo.Context) error {
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	resp, err := h.GetTopicsUseCase.Execute(c.Request().Context(), restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Get topic error")
		return response.Error(c, http.StatusBadRequest, "Failed to get topics")
	}
	return response.Success(c, resp, "Get topics successfully")
}

// DeleteTopicHandler godoc
// @Summary Delete topic
// @Description Delete topic by ID
// @Tags Topic
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Success 200 {object} app.DeleteTopicSuccessResponseDoc "Delete topic successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/topic/{id} [delete]
func (h *TopicHandler) Delete(c echo.Context) error {
	idInt64, err := parseRequiredIDParam(c.Param("id"), "missing topic id", "invalid topic id format")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	if err := h.DeleteUseCase.Execute(c.Request().Context(), idInt64, restaurantID); err != nil {
		h.Logger.Error().Err(err).Msg("Delete topic error")
		return response.Error(c, http.StatusBadRequest, "Failed to delete topic")
	}
	return response.Success[any](c, nil, "Delete topic successfully")
}
