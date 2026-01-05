package menuhttp

import (
	"go-ai/internal/identity/domain/rbac"
	"go-ai/internal/identity/transport/middlewares"

	"github.com/labstack/echo/v4"
)

const (
	optionItemIDPath  = "/option-item/:id"
	optionGroupIDPath = "/option-group/:id"
	topicIDPath       = "/topics/:id"
)

func RegisterMenuItemRoutes(api *echo.Group, h *MenuItemHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service) {
	r := api.Group("/menu", auth.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("/items", h.Create)
	r.GET("/items/:id", h.Get)
	r.PATCH("/items/:id", h.Update)
	r.POST("/items/search", h.Search)
	r.PATCH("/items/:id/status", h.UpdateStatus)
	r.DELETE("/items/:id", h.Delete)
}

func RegisterTopicRoutes(api *echo.Group, h *TopicHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service) {
	r := api.Group("/menu", auth.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("/topics", h.Create)
	r.GET(topicIDPath, h.Get)
	r.PUT(topicIDPath, h.Update)
	r.DELETE(topicIDPath, h.Delete)
	r.GET("/restaurant/topics", h.GetTopics)
	r.GET("/topics/combobox", h.GetCombobox)
}

func RegisterOptionGroupRoutes(api *echo.Group, h *OptionGroupHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service) {
	r := api.Group("/menu", auth.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("/option-group", h.Create)
	r.GET(optionGroupIDPath, h.Get)
	r.GET("/item/:id/option-groups", h.GetByMenuItem)
	r.PUT(optionGroupIDPath, h.Update)
	r.DELETE(optionGroupIDPath, h.Delete)
}

func RegisterOptionItemRoutes(api *echo.Group, h *OptionItemHandler, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service) {
	r := api.Group("/menu", auth.Handler, middlewares.RequirePermission(rbacSvc, "admin"))
	r.POST("/option-item", h.Create)
	r.GET(optionItemIDPath, h.Get)
	r.GET("/option-group/:id/option-items", h.GetByGroup)
	r.PUT(optionItemIDPath, h.Update)
	r.DELETE(optionItemIDPath, h.Delete)
}
