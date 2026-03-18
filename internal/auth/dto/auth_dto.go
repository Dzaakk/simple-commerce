package dto

type RegisterCustomerRequest struct {
	Email    string
	Password string
	FullName string
	Phone    string
}

type RegisterSellerRequest struct {
	Email    string
	Password string
	FullName string
	Phone    string
	ShopName string
}

type LoginRequest struct {
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required"`
	UserType string `json:"user_type" binding:"required,oneof=customer seller"`
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string
	ExpiresIn   int
}
