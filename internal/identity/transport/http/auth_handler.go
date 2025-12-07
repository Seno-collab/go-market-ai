package identityhttp

import (
	authapp "go-ai/internal/identity/application/auth"
	"go-ai/internal/transport/response"
	domainerr "go-ai/pkg/domain_err"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	RegisterUC       *authapp.RegisterUseCase
	LoginUC          *authapp.LoginUseCase
	RefreshTokenUC   *authapp.RefreshTokenUseCase
	ProfileUC        *authapp.GetProfileUseCase
	ChangePasswordUC *authapp.ChangePasswordUseCase
	LogoutUC         *authapp.LogoutUseCase
	Logger           zerolog.Logger
}

func NewAuthHandler(
	regUC *authapp.RegisterUseCase,
	loginUC *authapp.LoginUseCase,
	refreshUC *authapp.RefreshTokenUseCase,
	profileUC *authapp.GetProfileUseCase,
	changePasswordUC *authapp.ChangePasswordUseCase,
	logoutUC *authapp.LogoutUseCase,
	logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		RegisterUC:       regUC,
		LoginUC:          loginUC,
		RefreshTokenUC:   refreshUC,
		ProfileUC:        profileUC,
		ChangePasswordUC: changePasswordUC,
		LogoutUC:         logoutUC,
		Logger:           logger.With().Str("component", "Auth handler").Logger(),
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
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
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
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
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
		h.Logger.Error().Err(err).Msg("Bind json")
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}

	responseData, err := h.RefreshTokenUC.Execute(c.Request().Context(), in)
	if err != nil {
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
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
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	resp := &authapp.GetProfileResponse{
		Email:    profile.Email,
		FullName: profile.FullName,
		Role:     profile.Role,
		IsActive: profile.IsActive,
	}
	return response.Success[authapp.GetProfileResponse](c, resp, "Profile retrieved successfully")
}

// ChangePassword godoc
// @Summary Change user password
// @Description Allows the authenticated user to change their password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authapp.ChangePasswordRequest true "Old and new password"
// @Success 200 {object} app.ChangePasswordSuccessResponseDoc "Password changed successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/auth/change-password [patch]
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	var in authapp.ChangePasswordRequest
	if err := c.Bind(&in); err != nil {
		h.Logger.Error().Err(err).Msg("Bind json")
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	userID := c.Get("user_id")
	if userID == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	err := h.ChangePasswordUC.Execute(c.Request().Context(), in, userUUID)
	if err != nil {
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "internal server error")
	}
	return response.Success[any](c, nil, "password changed successfully")
}

// Logout godoc
// @Summary Logout user
// @Description Invalidates the current user's authentication token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} app.LogoutSuccessResponseDoc "Successfully logged out"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	err := h.LogoutUC.Execute(c.Request().Context(), userUUID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Delete cache error")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success[any](c, nil, "Log out successfully")
}
