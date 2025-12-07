package response

import "errors"

var (
	ErrInvalidEmail            = errors.New("invalid email")
	ErrInvalidField            = errors.New("Invalid field")
	ErrNotFound                = errors.New("user not found")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrUserInactive            = errors.New("user is inactive")
	ErrTokenInvalid            = errors.New("Token is invalid")
	ErrTokenExpired            = errors.New("Token is expired")
	ErrTokenMissing            = errors.New("Token is missing")
	ErrTokenMalformed          = errors.New("Token is malformed")
	ErrNameAlreadyExists       = errors.New("name already exists")
	ErrUnauthorizedAccess      = errors.New("unauthorized access")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrConflict                = errors.New("email already exists")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrInvalidPhoneNumber      = errors.New("Invalid phone number")
	ErrInternalServerError     = errors.New("internal server error")
	ErrRestaurantNoExitis      = errors.New("Restaurant not exitis")
	ErrRestaurantNameExitis    = errors.New("Name restaurant exitis")
	ErrTokenNotActive          = errors.New("Token is not active yet")
	ErrTokenGenerateFail       = errors.New("Failed to generate token")
	ErrorRefreshTokenEmpty     = errors.New("Refresh Token empty string")
	ErrPasswordVerifyFail      = errors.New("Password verification failed")
	ErrTokenWrongSigningMethod = errors.New("Token has wrong signing method")
)
