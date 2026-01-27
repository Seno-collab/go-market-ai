package menuhttp

import (
	topicapp "go-ai/internal/menu/application/topic"
	"go-ai/internal/transport/response"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
)

type TopicHandler struct {
	CreateUseCase      *topicapp.CreateUseCase
	GetUseCase         *topicapp.GetUseCase
	GetTopicsUseCase   *topicapp.GetTopicsUseCase
	UpdateUseCase      *topicapp.UpdateUseCase
	DeleteUseCase      *topicapp.DeleteUseCase
	GetComboboxUseCase *topicapp.ComboboxUseCase
	Logger             zerolog.Logger
}

func NewTopicHandler(
	createUseCase *topicapp.CreateUseCase,
	getUseCase *topicapp.GetUseCase,
	getTopicsUseCase *topicapp.GetTopicsUseCase,
	updateUseCase *topicapp.UpdateUseCase,
	deleteUseCase *topicapp.DeleteUseCase,
	getComboboxUseCase *topicapp.ComboboxUseCase,
	logger zerolog.Logger,
) *TopicHandler {
	return &TopicHandler{
		CreateUseCase:      createUseCase,
		GetUseCase:         getUseCase,
		GetTopicsUseCase:   getTopicsUseCase,
		UpdateUseCase:      updateUseCase,
		DeleteUseCase:      deleteUseCase,
		GetComboboxUseCase: getComboboxUseCase,
		Logger:             logger.With().Str("handler", "TopicHandler").Logger(),
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
// @Router       /api/menu/topics [post]
func (h *TopicHandler) Create(c *echo.Context) error {
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
// @Router /api/menu/topics/{id} [get]
func (h *TopicHandler) Get(c *echo.Context) error {
	idInt64, err := parseRequiredIDParam(c.Param("id"), "Missing topic ID", "Invalid topic ID format")
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
// @Router /api/menu/topics/{id} [put]
func (h *TopicHandler) Update(c *echo.Context) error {
	idInt64, err := parseRequiredIDParam(c.Param("id"), "Missing topic ID", "Invalid topic ID format")
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
// @Param name query string false "Topic name keyword"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} app.GetTopicsByRestaurantSuccessResponseDoc "Get topics by restaurant successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/topics/search [get]
func (h *TopicHandler) GetTopics(c *echo.Context) error {
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	page := 1
	limit := 20
	name := strings.TrimSpace(c.QueryParam("name"))
	if v := c.QueryParam("page"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil || p < 1 || p > math.MaxInt32 {
			return response.Error(c, http.StatusBadRequest, "invalid page")
		}
		page = p
	}
	if v := c.QueryParam("limit"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil || l < 1 || l > 100 || l > math.MaxInt32 {
			return response.Error(c, http.StatusBadRequest, "invalid limit")
		}
		limit = l
	}
	resp, err := h.GetTopicsUseCase.Execute(c.Request().Context(), restaurantID, name, int32(page), int32(limit))
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
// @Router /api/menu/topics/{id} [delete]
func (h *TopicHandler) Delete(c *echo.Context) error {
	idInt64, err := parseRequiredIDParam(c.Param("id"), "Missing topic ID", "Invalid topic ID format")
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

// GetTopicComboboxHandler godoc
// @Summary Get topic combobox
// @Description Get list of topics for combobox by restaurant (optional parent_id)
// @Tags Topic
// @Accept json
// @Produce json
// @Success 200 {array} app.GetTopicsByRestaurantComboboxSuccessResponseDoc "Get topic combobox successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/menu/topics/combobox [get]
func (h *TopicHandler) GetCombobox(c *echo.Context) error {
	restaurantID, err := getRestaurantID(c)
	if err != nil {
		h.Logger.Error().Err(err).Msg(logInvalidRestaurantID)
		return response.Error(c, http.StatusBadRequest, errInvalidRestaurantID)
	}
	resp, err := h.GetComboboxUseCase.Execute(c.Request().Context(), restaurantID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Get topic combobox error")
		return response.Error(c, http.StatusBadRequest, "Failed to get topic combobox")
	}
	return response.Success(c, resp, "Get topic combobox successfully")
}
