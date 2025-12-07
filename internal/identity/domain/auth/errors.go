package auth

import (
	domainerr "go-ai/pkg/domain_err"
	"net/http"
)

var (
	ErrFullNameRequired = domainerr.New(http.StatusBadRequest, "full name is required")
	ErrPasswordTooShort = domainerr.New(http.StatusBadRequest, "password must be at least 6 characters")
	ErrRoleRequired     = domainerr.New(http.StatusBadRequest, "role is required")
	// 400 — Validation
	ErrInvalidEmail       = domainerr.New(http.StatusBadRequest, "invalid email")
	ErrWeakPassword       = domainerr.New(http.StatusBadRequest, "password must contain uppercase, lowercase, digit and special character")
	ErrInvalidPhoneNumber = domainerr.New(http.StatusBadRequest, "invalid phone number")

	ErrUnauthorizedAccess = domainerr.New(http.StatusUnauthorized, "unauthorized access")

	ErrUserNotFound = domainerr.New(http.StatusNotFound, "user not found")

	ErrNameAlreadyExists    = domainerr.New(http.StatusConflict, "name already exists")
	ErrUserAlreadyExists    = domainerr.New(http.StatusConflict, "user already exists")
	ErrEmailAlreadyExists   = domainerr.New(http.StatusConflict, "email already exists")
	ErrRestaurantNameExists = domainerr.New(http.StatusConflict, "restaurant name already exists")

	// Token errors → 401 (Unauthorized)
	ErrTokenInvalid            = domainerr.New(http.StatusUnauthorized, "token is invalid")
	ErrTokenExpired            = domainerr.New(http.StatusUnauthorized, "token is expired")
	ErrTokenMissing            = domainerr.New(http.StatusUnauthorized, "token is missing")
	ErrTokenMalformed          = domainerr.New(http.StatusUnauthorized, "token is malformed")
	ErrTokenNotActive          = domainerr.New(http.StatusUnauthorized, "token is not active yet")
	ErrTokenGenerateFail       = domainerr.New(http.StatusInternalServerError, "failed to generate token")
	ErrTokenWrongSigningMethod = domainerr.New(http.StatusUnauthorized, "token has wrong signing method")

	// Refresh token errors
	ErrorRefreshTokenEmpty = domainerr.New(http.StatusBadRequest, "refresh token is empty")

	// Password verify
	ErrPasswordVerifyFail = domainerr.New(http.StatusBadRequest, "password verification failed")

	ErrUserInactive = domainerr.New(http.StatusForbidden, "user is inactive")
)
