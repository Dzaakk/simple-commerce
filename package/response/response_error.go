package response

import "errors"

var (
	// Auth
	ErrEmailAlreadyExist  = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("email or password is incorrect")
	ErrEmailNotVerified   = errors.New("email not verified, please check your inbox")

	// Activation
	ErrInvalidActivationCode = errors.New("invalid or expired activation code")

	// Token
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

var (
	ErrUserNotFound = errors.New("user not found")
)
