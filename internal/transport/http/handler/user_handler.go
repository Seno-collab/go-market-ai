package handler

import (
	"go-ai/internal/domain/user"
	"go-ai/pkg/common"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	svc      *user.Service
	validate *validator.Validate
}

func NewUserHandler(svc *user.Service, v *validator.Validate) *UserHandler {
	return &UserHandler{
		svc:      svc,
		validate: v,
	}
}

type registerReq struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name"`
	Password string `json:"password"  validate:"required,password"`
}
type userResp struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *UserHandler) Register(c echo.Context) error {
	var in registerReq
	if err := c.Bind(&in); err != nil {
		return common.ErrorResponse(
			c, http.StatusBadRequest, "Invalid JSON body",
			[]common.ErrorDetail{{Field: "body", Message: err.Error()}},
		)
	}
	if h.validate != nil {
		if err := h.validate.Struct(in); err != nil {
			// return common.ErrorResponse(c, http.StatusBadRequest, "Validation failed", common.MapValidationErrors(err))
		}
	}
	id, err := h.svc.Register(user.User{Email: in.Email, Name: in.Name})
	if err != nil {
		switch err {
		case user.ErrInvalidEmail, user.ErrInvalidName:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case user.ErrConflict:
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal"})
		}
	}
	return c.JSON(http.StatusCreated, map[string]any{"id": id})
}

func (h *UserHandler) Detail(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	u, err := h.svc.Detail(id)
	if err != nil {
		if err == user.ErrNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal"})
	}
	return c.JSON(http.StatusOK, userResp{ID: u.ID, Email: u.Email, Name: u.Name})
}
