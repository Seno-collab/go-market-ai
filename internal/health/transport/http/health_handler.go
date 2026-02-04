package healthhttp

import (
	healthapp "go-ai/internal/health/application/health"
	"go-ai/pkg/response"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
)

type HealthHandler struct {
	checker *healthapp.CheckHealthUseCase
	logger  zerolog.Logger
}

func NewHealthHandler(checker *healthapp.CheckHealthUseCase, logger zerolog.Logger) *HealthHandler {
	return &HealthHandler{
		checker: checker,
		logger:  logger.With().Str("component", "HealthHandler").Logger(),
	}
}

// Health godoc
// @Summary Health check
// @Description Provides service health along with dependency status.
// @Tags Health
// @Produce json
// @Success 200 {object} app.HealthSuccessResponseDoc "Service healthy"
// @Failure 503 {object} app.HealthFailureResponseDoc "Service degraded or down"
// @Router /api/health [get]
func (h *HealthHandler) Health(c *echo.Context) error {
	result, ok := h.checker.Execute(c.Request().Context())
	if !ok || result.Status != healthapp.StatusUp {
		h.logger.Warn().Any("health", result).Msg("health check degraded")
		return response.Error(c, http.StatusServiceUnavailable, "Service degraded")
	}
	return response.Success(c, result, "Service healthy")
}
