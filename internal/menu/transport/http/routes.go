package menuhttp

import (
	"go-ai/internal/identity/domain/rbac"
	"go-ai/internal/identity/transport/middlewares"
	middlewaresRestaurant "go-ai/internal/restaurant/transport/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterMenuRoutes(api *echo.Group, h *MenuItemHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service, requiredRestaurant *middlewaresRestaurant.RestaurantRequired) {
	r := api.Group("/menu", auth.Handler, requiredRestaurant.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("/item", h.Create)
}
