package model

import "database/sql"

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CustomerRegistrationReq struct {
	Username    string       `json:"username" validate:"required"`
	Email       string       `json:"email" validate:"required,email"`
	PhoneNumber string       `json:"phone_number" validate:"required"`
	Password    string       `json:"password" validate:"required"`
	Gender      int          `json:"gender"`
	DateOfBirth sql.NullTime `json:"date_of_birth"`
}

type SellerRegistrationReq struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CustomerActivationReq struct {
	Email          string `json:"email" validate:"required,email"`
	ActivationCode int    `json:"activation_code"`
}
