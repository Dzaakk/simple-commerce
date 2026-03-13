package constant

// UserType represents the type of user in the system.
type UserType string

const (
	Customer UserType = "customer"
	Seller   UserType = "seller"
)

// ActivationCodeType represents the type of activation code.
type ActivationCodeType string

const (
	EmailVerification ActivationCodeType = "email_verification"
	ResetPassword     ActivationCodeType = "password_reset"
)

// UserStatus represents the status of a user account.
type UserStatus string

const (
	StatusPending UserStatus = "pending"
	StatusActive  UserStatus = "active"
)
