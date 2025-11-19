package handler

import (
	domain "go-ai/internal/domain/auth"
	"go-ai/internal/domain/user"
	authservice "go-ai/internal/service/auth"
	"go-ai/pkg/common"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	svc *authservice.Service
}

func NewAuthHandler(svc *authservice.Service) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

type ErrorResponseDoc struct {
	ResponseCode string              `json:"response_code,omitempty"`
	Message      string              `json:"message"`
	Error        *common.ErrorDetail `json:"error,omitempty"`
}
type RegisterReq struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name"`
	Password string `json:"password"  validate:"required,password"`
}
type RegisterSuccess struct {
}

type RegisterSuccessResponse struct {
	Message      string           `json:"message"`
	ResponseCode string           `json:"response_code,omitempty"`
	Data         *RegisterSuccess `json:"data,omitempty"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and full name
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterReq true "User registration payload"
// @Success 200 {object} RegisterSuccessResponse "User created successfully"
// @BadRequest 400 {object} ErrorResponseDoc "Invalid input or validation failed"
// @Conflict 409 {object} ErrorResponseDoc "User already exists"
// @InternalServerError 500 {object} ErrorResponseDoc "Internal server error"
// @Router /users/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var in RegisterReq
	if err := c.Bind(&in); err != nil {

	}
	// if h.validate != nil {
	// 	if err := h.validate.Struct(in); err != nil {
	// 		return common.ErrorResponse(c, http.StatusBadRequest, "Validation failed")
	// 	}
	// }
	_, err := h.svc.Register(&domain.Auth{Email: in.Email, FullName: in.FullName})
	if err != nil {
		switch err {
		case user.ErrInvalidEmail, user.ErrInvalidName:
			return common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		case user.ErrConflict:
			return common.ErrorResponse(c, http.StatusConflict, err.Error())
		default:
			return common.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	return common.SuccessResponse[any](c, nil, "create user success")
}
