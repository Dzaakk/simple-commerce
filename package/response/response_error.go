package response

import "errors"

var (
	// Auth
	ErrEmailAlreadyExist  = errors.New("unable to complete registration. please try again or contact support")
	ErrInvalidCredentials = errors.New("email or password is incorrect")
	ErrEmailNotVerified   = errors.New("email is not verified")

	// Activation
	ErrInvalidActivationCode = errors.New("activation code is invalid")
	// ErrExpiredCode     = errors.New("activation code has expired")
	// ErrCodeAlreadyUsed = errors.New("activation code has already been used")

	// Token
	// ErrInvalidToken = errors.New("token is invalid")
	// ErrExpiredToken = errors.New("token has expired")
	// ErrRevokedToken = errors.New("token has been revoked")
)

var (
	// ErrInternalServer = errors.New("internal server error")
	// ErrBadRequest     = errors.New("bad request")
	// ErrUnauthorized   = errors.New("unauthorized")
	// ErrForbidden      = errors.New("forbidden")
	ErrUserNotFound = errors.New("user not found")
)
