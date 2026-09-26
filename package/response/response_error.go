package response

import "errors"

var (
	ErrEmailAlreadyExist   = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("email or password is incorrect")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrUserNotFound        = errors.New("user not found")
)
