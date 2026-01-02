package auth

import (
	domainerr "go-ai/pkg/domain_err"
	"net/http"
)

var (
	ErrUserInactive            = domainerr.New(http.StatusForbidden, "User is inactive")
	ErrRoleRequired            = domainerr.New(http.StatusBadRequest, "Role is required")
	ErrUserNotFound            = domainerr.New(http.StatusNotFound, "User not found")
	ErrInvalidEmail            = domainerr.New(http.StatusBadRequest, "Invalid email")
	ErrFullNameRequired        = domainerr.New(http.StatusBadRequest, "Full name is required")
	ErrNameAlreadyExists       = domainerr.New(http.StatusConflict, "Name already exists")
	ErrUserAlreadyExists       = domainerr.New(http.StatusConflict, "User already exists")
	ErrEmailAlreadyExists      = domainerr.New(http.StatusConflict, "Email already exists")
	ErrInvalidPhoneNumber      = domainerr.New(http.StatusBadRequest, "Invalid phone number")
	ErrUnauthorizedAccess      = domainerr.New(http.StatusUnauthorized, "Unauthorized access")
	ErrTokenInvalid            = domainerr.New(http.StatusUnauthorized, "Token is invalid")
	ErrTokenExpired            = domainerr.New(http.StatusUnauthorized, "Token is expired")
	ErrTokenMissing            = domainerr.New(http.StatusUnauthorized, "Token is missing")
	ErrorRefreshTokenEmpty     = domainerr.New(http.StatusBadRequest, "Refresh token is empty")
	ErrTokenMalformed          = domainerr.New(http.StatusUnauthorized, "Token is malformed")
	ErrOldPasswordIncorrect    = domainerr.New(http.StatusBadRequest, "Old password is incorrect")
	ErrPasswordVerifyFail      = domainerr.New(http.StatusBadRequest, "Password verification failed")
	ErrRestaurantNameExists    = domainerr.New(http.StatusConflict, "Restaurant name already exists")
	ErrTokenNotActive          = domainerr.New(http.StatusUnauthorized, "Token is not active yet")
	ErrInvalidCredentials      = domainerr.New(http.StatusBadRequest, "Email or password is incorrect")
	ErrHashPasswordFailed      = domainerr.New(http.StatusInternalServerError, "Failed to hash password")
	ErrPasswordTooShort        = domainerr.New(http.StatusBadRequest, "Password must be at least 6 characters")
	ErrTokenWrongSigningMethod = domainerr.New(http.StatusUnauthorized, "Token has wrong signing method")
	ErrTokenGenerateFail       = domainerr.New(http.StatusInternalServerError, "Failed to generate token")
	ErrConfirmPassword         = domainerr.New(http.StatusBadRequest, "New password and confirm password do not match")
	ErrWeakPassword            = domainerr.New(http.StatusBadRequest, "Password must contain uppercase, lowercase, digit and special character")
)
