package response

import "errors"

var (
	// Auth
	ErrEmailAlreadyExist = errors.New("unable to complete registration. please try again or contact support")
	// ErrInvalidCredentials = errors.New("email or password is incorrect")
	// ErrAccountNotActive   = errors.New("account is not active, please verify your email")

	// Activation
	// ErrInvalidCode     = errors.New("activation code is invalid")
	// ErrExpiredCode     = errors.New("activation code has expired")
	// ErrCodeAlreadyUsed = errors.New("activation code has already been used")

	// Token
	// ErrInvalidToken = errors.New("token is invalid")
	// ErrExpiredToken = errors.New("token has expired")
	// ErrRevokedToken = errors.New("token has been revoked")
)
