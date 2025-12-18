package menuhttp

import (
	"go-ai/internal/identity/domain/rbac"
	"go-ai/internal/identity/transport/middlewares"
	middlewaresRestaurant "go-ai/internal/restaurant/transport/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterMenuItemRoutes(api *echo.Group, h *MenuItemHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service, requiredRestaurant *middlewaresRestaurant.RestaurantRequired) {
	r := api.Group("/menu", auth.Handler, requiredRestaurant.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("/item", h.Create)
	r.GET("/item/:id", h.Get)
	r.PUT("item/:id", h.Update)
	r.GET("/restaurant/items", h.GetItems)
}

func RegisterTopicRoutes(api *echo.Group, h *TopicHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service, requiredRestaurant *middlewaresRestaurant.RestaurantRequired) {
	r := api.Group("/menu", auth.Handler, requiredRestaurant.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("/topic", h.Create)
	r.GET("/topic/:id", h.Get)
	r.GET("/restaurant/topics", h.GetTopcis)
}
