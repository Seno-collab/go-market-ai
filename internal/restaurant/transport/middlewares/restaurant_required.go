package middlewares

// import (
// 	"go-ai/internal/identity/domain/rbac"
// 	repo "go-ai/internal/restaurant/infrastructure/db"
// 	"go-ai/internal/transport/response"
// 	"net/http"

// 	"github.com/google/uuid"
// 	"github.com/labstack/echo/v4"
// )

// type RestaurantRequired struct {
// 	queries *repo.RestaurantRepo
// }

// func New(queries *repo.RestaurantRepo) *RestaurantRequired {
// 	return &RestaurantRequired{
// 		queries: queries,
// 	}
// }

// func (r *RestaurantRequired) Handler(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		userID := c.Get("user_id")
// 		role := c.Get("role")
// 		if userID == nil {
// 			return response.Error(c, 401, "Unauthorized access")
// 		}
// 		userUUID, ok := userID.(uuid.UUID)
// 		if !ok {
// 			return response.Error(c, http.StatusForbidden, "You do not have permission to access this resource")
// 		}
// 		if role != rbac.Admin {
// 			id, err := r.queries.GetRestaurantByUserID(c.Request().Context(), userUUID)
// 			if err != nil {
// 				return response.Error(c, http.StatusForbidden, "You do not have permission to access this resource")
// 			}
// 			c.Set("restaurant_id", id)
// 		}
// 		return next(c)
// 	}
// }
