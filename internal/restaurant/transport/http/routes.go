package restauranthttp

import (
	"go-ai/internal/identity/domain/rbac"
	middlewares "go-ai/internal/identity/transport/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterRestaurantRoutes(api *echo.Group, h *RestaurantHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service) {
	r := api.Group("/restaurants", auth.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("", h.Create)
	r.GET("/:id", h.GetByID)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}
