package model

import (
	"database/sql"
	"time"
)

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type CustomerToken struct {
	Email string `json:"email"`
}

type CustomerRegistration struct {
	Username    string       `json:"username" validate:"required"`
	Email       string       `json:"email" validate:"required,email"`
	PhoneNumber string       `json:"phone_number" validate:"required"`
	Password    string       `json:"password" validate:"required"`
	Gender      int          `json:"gender" validate:"number"`
	DateOfBirth sql.NullTime `json:"date_of_birth"`
}

type SellerRegistration struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TCustomerActivationCode struct {
	CustomerID     int64        `json:"customer_id"`
	CodeActivation string       `json:"code_activation"`
	IsUsed         bool         `json:"is_used"`
	CreatedAt      time.Time    `json:"created_at"`
	UsedAt         sql.NullTime `json:"used_at"`
}
type CustomerActivationCode struct {
	CustomerID     string `json:"customer_id"`
	CodeActivation string `json:"code_activation"`
	IsUsed         string `json:"is_used"`
}
type TSellerActivationCode struct {
	SellerID       int64        `json:"customer_id"`
	CodeActivation string       `json:"code_activation"`
	IsUsed         bool         `json:"is_used"`
	CreatedAt      time.Time    `json:"created_at"`
	UsedAt         sql.NullTime `json:"used_at"`
}
type SellerActivationCode struct {
	SellerID       string `json:"customer_id"`
	CodeActivation string `json:"code_activation"`
	IsUsed         string `json:"is_used"`
}
