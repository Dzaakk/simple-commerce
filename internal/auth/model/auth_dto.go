package model

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CustomerRegistrationReq struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

type SellerRegistrationReq struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	StoreName   string `json:"store_name" validate:"required"`
	Address     string `json:"address" validate:"required"`
}

type ActivationReq struct {
	Email          string `json:"email" validate:"required,email"`
	ActivationCode string `json:"activation_code"`
}
