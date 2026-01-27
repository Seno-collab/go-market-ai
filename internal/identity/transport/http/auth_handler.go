package identityhttp

import (
	authapp "go-ai/internal/identity/application/auth"
	"go-ai/internal/transport/response"
	domainerr "go-ai/pkg/domain_err"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	RegisterUseCase       *authapp.RegisterUseCase
	LoginUseCase          *authapp.LoginUseCase
	RefreshTokenUseCase   *authapp.RefreshTokenUseCase
	ProfileUseCase        *authapp.GetProfileUseCase
	ChangePasswordUseCase *authapp.ChangePasswordUseCase
	UpdateProfileUseCase  *authapp.UpdateProfileUseCase
	LogoutUseCase         *authapp.LogoutUseCase
	Logger                zerolog.Logger
}

func NewAuthHandler(
	registerUseCase *authapp.RegisterUseCase,
	loginUseCase *authapp.LoginUseCase,
	refreshTokenUseCase *authapp.RefreshTokenUseCase,
	profileUseCase *authapp.GetProfileUseCase,
	changePasswordUseCase *authapp.ChangePasswordUseCase,
	updateProfileUseCase *authapp.UpdateProfileUseCase,
	logoutUseCase *authapp.LogoutUseCase,
	logger zerolog.Logger,
) *AuthHandler {
	return &AuthHandler{
		RegisterUseCase:       registerUseCase,
		LoginUseCase:          loginUseCase,
		RefreshTokenUseCase:   refreshTokenUseCase,
		ProfileUseCase:        profileUseCase,
		ChangePasswordUseCase: changePasswordUseCase,
		UpdateProfileUseCase:  updateProfileUseCase,
		LogoutUseCase:         logoutUseCase,
		Logger:                logger.With().Str("component", "AuthHandler").Logger(),
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
func (h *AuthHandler) Register(c *echo.Context) error {
	var in authapp.RegisterRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	_, err := h.RegisterUseCase.Execute(c.Request().Context(), in)
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
func (h *AuthHandler) Login(c *echo.Context) error {
	var in authapp.LoginRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	responseData, err := h.LoginUseCase.Execute(c.Request().Context(), in)
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
	return response.Success(c, responseData, "Login success")
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
func (h *AuthHandler) RefreshToken(c *echo.Context) error {
	var in authapp.RefreshTokenRequest
	if err := c.Bind(&in); err != nil {
		h.Logger.Error().Err(err).Msg("Bind json")
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}

	responseData, err := h.RefreshTokenUseCase.Execute(c.Request().Context(), in)
	if err != nil {
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success(c, responseData, "Token refreshed successfully")
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
func (h *AuthHandler) GetProfile(c *echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	profile, err := h.ProfileUseCase.Execute(c.Request().Context(), userUUID)
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
		ImageUrl: profile.ImageUrl,
	}
	return response.Success(c, resp, "Profile retrieved successfully")
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile information of the authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authapp.UpdateProfileRequest true "Profile update payload"
// @Success 200 {object} app.UpdateProfileSuccessResponseDoc "Profile updated successfully"
// @Failure default {object} response.ErrorDoc "Errors"
// @Router /api/auth/profile [patch]
func (h *AuthHandler) UpdateProfile(c *echo.Context) error {
	var in authapp.UpdateProfileRequest
	if err := c.Bind(&in); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}
	userID := c.Get("user_id")
	if userID == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to update profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	profile, err := h.UpdateProfileUseCase.Execute(c.Request().Context(), userUUID, in)
	if err != nil {
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success(c, profile, "Profile updated successfully")
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
func (h *AuthHandler) ChangePassword(c *echo.Context) error {
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
	err := h.ChangePasswordUseCase.Execute(c.Request().Context(), in, userUUID)
	if err != nil {
		if ae, ok := err.(domainerr.AppError); ok {
			return response.Error(c, ae.Status, ae.Msg)
		}
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success[any](c, nil, "Password changed successfully")
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
func (h *AuthHandler) Logout(c *echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.Logger.Error().Msg("failed to get profile: invalid user ID type")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	err := h.LogoutUseCase.Execute(c.Request().Context(), userUUID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Delete cache error")
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}
	return response.Success[any](c, nil, "Log out successfully")
}
