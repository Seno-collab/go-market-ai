package healthhttp

import "github.com/labstack/echo/v5"

func RegisterHealthRoutes(api *echo.Group, h *HealthHandler) {
	api.GET("/health", h.Health)
}
