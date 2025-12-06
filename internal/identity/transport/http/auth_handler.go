package identityhttp

import (
	"errors"

	authapp "go-ai/internal/identity/application/auth"
	"go-ai/internal/transport/response"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	RegisterUC     *authapp.RegisterUseCase
	LoginUC        *authapp.LoginUseCase
	RefreshTokenUC *authapp.RefreshTokenUseCase
	ProfileUC      *authapp.GetProfileUseCase
	Logger         zerolog.Logger
}

func NewAuthHandler(
	regUC *authapp.RegisterUseCase,
	loginUC *authapp.LoginUseCase,
	refreshUC *authapp.RefreshTokenUseCase,
	profileUC *authapp.GetProfileUseCase,
	logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		RegisterUC:     regUC,
		LoginUC:        loginUC,
		RefreshTokenUC: refreshUC,
		ProfileUC:      profileUC,
		Logger:         logger.With().Str("component", "Auth handler").Logger(),
	}
}

// Register godoc
// @Summary Register a new userRegisterRequest
// @Description Create a new user account with email and full name
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authapp.RegisterRequest true "User registration payload"
// @Success 200 {object} app.RegisterSuccessResponseDoc "User created successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var in authapp.RegisterRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	_, err := h.RegisterUC.Execute(c.Request().Context(), in)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to register user")
		switch err {
		case response.ErrInvalidEmail, response.ErrInvalidName, response.ErrInvalidPassword, response.ErrUserAlreadyExists, response.ErrNameAlreadyExists:
			details := response.ErrorDetail{}
			if errors.Is(err, response.ErrInvalidEmail) {
				details = response.ErrorDetail{
					Field:   "email",
					Message: "Email is a required field",
				}
			}
			if errors.Is(err, response.ErrInvalidPassword) {
				details = response.ErrorDetail{
					Field:   "password",
					Message: "Password is a required field",
				}
			}
			if errors.Is(err, response.ErrInvalidName) {
				details = response.ErrorDetail{
					Field:   "name",
					Message: "Name is a required field",
				}
			}
			return response.Error(c, http.StatusBadRequest, err.Error(), details)
		case response.ErrConflict:
			return response.Error(c, http.StatusConflict, err.Error())
		default:
			return response.Error(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	return response.Success[any](c, nil, "Create user success")
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authapp.LoginRequest true "User login payload"
// @Success 200 {object} app.LoginSuccessResponseDoc "Login successful"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var in authapp.LoginRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	responseData, err := h.LoginUC.Execute(c.Request().Context(), in)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to login user")
		switch err {
		case response.ErrInvalidEmail, response.ErrInvalidPassword, response.ErrNotFound, response.ErrPasswordVerifyFail, response.ErrUserInactive:
			details := response.ErrorDetail{}
			if errors.Is(err, response.ErrInvalidEmail) {
				details = response.ErrorDetail{
					Field:   "email",
					Message: "Email is a required field",
				}
			}
			if errors.Is(err, response.ErrInvalidPassword) {
				details = response.ErrorDetail{
					Field:   "password",
					Message: "Password is a required field",
				}
			}
			return response.Error(c, http.StatusBadRequest, "Invalid email or password", details)
		default:
			return response.Error(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	if responseData.AccessToken == "" || responseData.RefreshToken == "" {
		h.Logger.Error().Msg("Failed to login user: invalid credentials")
		return response.Error(c, http.StatusBadRequest, "Invalid email or password")
	}
	return response.Success[authapp.LoginResponse](c, responseData, "login success")
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate a new access token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authapp.RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} app.RefreshTokenSuccessResponseDoc "Token refreshed successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var in authapp.RefreshTokenRequest
	if err := c.Bind(&in); err != nil {
		h.Logger.Error().Err(err).Msg("")
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}

	responseData, err := h.RefreshTokenUC.Execute(c.Request().Context(), in)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to refresh token")
		switch err {
		case response.ErrorRefreshTokenEmpty:
			details := response.ErrorDetail{}
			if errors.Is(err, response.ErrorRefreshTokenEmpty) {
				details = response.ErrorDetail{
					Field:   "refresh_token",
					Message: "Refresh token is a required field",
				}
			}
			return response.Error(c, http.StatusBadRequest, "Invalid refresh token", details)
		case response.ErrTokenInvalid, response.ErrTokenExpired, response.ErrTokenMalformed:
			return response.Error(c, http.StatusBadRequest, err.Error())
		default:
			return response.Error(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	return response.Success[authapp.RefreshTokenResponse](c, responseData, "Token refreshed successfully")
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve the profile information of the authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} app.GetProfileSuccessResponseDoc "Profile retrieved successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	profile, err := h.ProfileUC.Execute(c.Request().Context(), userUUID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to get profile")
		switch err {
		case response.ErrNotFound:
			return response.Error(c, http.StatusNotFound, "User not found")
		default:
			return response.Error(c, http.StatusInternalServerError, "Internal server error")
		}
	}
	resp := &authapp.GetProfileResponse{
		Email:    profile.Email,
		FullName: profile.FullName,
		Role:     profile.Role,
		IsActive: profile.IsActive,
	}
	return response.Success[authapp.GetProfileResponse](c, resp, "Profile retrieved successfully")
}
