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
	r.PUT("/topic/:id", h.Update)
	r.DELETE("/topic/:id", h.Delete)
	r.GET("/restaurant/topics", h.GetTopics)
}

func RegisterOptionGroupRoutes(api *echo.Group, h *OptionGroupHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service, requiredRestaurant *middlewaresRestaurant.RestaurantRequired) {
	r := api.Group("/menu", auth.Handler, requiredRestaurant.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("/option-group", h.Create)
	r.GET("/option-group/:id", h.Get)
	r.GET("/item/:id/option-groups", h.GetByMenuItem)
	r.PUT("/option-group/:id", h.Update)
	r.DELETE("/option-group/:id", h.Delete)
}

func RegisterOptionItemRoutes(api *echo.Group, h *OptionItemHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service, requiredRestaurant *middlewaresRestaurant.RestaurantRequired) {
	r := api.Group("/menu", auth.Handler, requiredRestaurant.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("/option-item", h.Create)
	r.GET("/option-item/:id", h.Get)
	r.GET("/option-group/:id/option-items", h.GetByGroup)
	r.PUT("/option-item/:id", h.Update)
	r.DELETE("/option-item/:id", h.Delete)
}
