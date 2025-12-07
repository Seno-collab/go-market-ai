package auth

import (
	domainerr "go-ai/pkg/domain_err"
	"net/http"
)

var (
	ErrUserInactive            = domainerr.New(http.StatusForbidden, "user is inactive")
	ErrRoleRequired            = domainerr.New(http.StatusBadRequest, "role is required")
	ErrUserNotFound            = domainerr.New(http.StatusNotFound, "user not found")
	ErrInvalidEmail            = domainerr.New(http.StatusBadRequest, "invalid email")
	ErrFullNameRequired        = domainerr.New(http.StatusBadRequest, "full name is required")
	ErrNameAlreadyExists       = domainerr.New(http.StatusConflict, "name already exists")
	ErrUserAlreadyExists       = domainerr.New(http.StatusConflict, "user already exists")
	ErrEmailAlreadyExists      = domainerr.New(http.StatusConflict, "email already exists")
	ErrInvalidPhoneNumber      = domainerr.New(http.StatusBadRequest, "invalid phone number")
	ErrUnauthorizedAccess      = domainerr.New(http.StatusUnauthorized, "unauthorized access")
	ErrTokenInvalid            = domainerr.New(http.StatusUnauthorized, "token is invalid")
	ErrTokenExpired            = domainerr.New(http.StatusUnauthorized, "token is expired")
	ErrTokenMissing            = domainerr.New(http.StatusUnauthorized, "token is missing")
	ErrorRefreshTokenEmpty     = domainerr.New(http.StatusBadRequest, "refresh token is empty")
	ErrTokenMalformed          = domainerr.New(http.StatusUnauthorized, "token is malformed")
	ErrOldPasswordIncorrect    = domainerr.New(http.StatusBadRequest, "old password is incorrect")
	ErrPasswordVerifyFail      = domainerr.New(http.StatusBadRequest, "password verification failed")
	ErrRestaurantNameExists    = domainerr.New(http.StatusConflict, "restaurant name already exists")
	ErrTokenNotActive          = domainerr.New(http.StatusUnauthorized, "token is not active yet")
	ErrInvalidCredentials      = domainerr.New(http.StatusBadRequest, "email or password is incorrect")
	ErrHashPasswordFailed      = domainerr.New(http.StatusInternalServerError, "failed to hash password")
	ErrPasswordTooShort        = domainerr.New(http.StatusBadRequest, "password must be at least 6 characters")
	ErrTokenWrongSigningMethod = domainerr.New(http.StatusUnauthorized, "token has wrong signing method")
	ErrTokenGenerateFail       = domainerr.New(http.StatusInternalServerError, "failed to generate token")
	ErrConfirmPassword         = domainerr.New(http.StatusBadRequest, "new password and confirm password do not match")
	ErrWeakPassword            = domainerr.New(http.StatusBadRequest, "password must contain uppercase, lowercase, digit and special character")
)
