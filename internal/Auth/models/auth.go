package models

import (
	"database/sql"
	"time"
)

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CustomerRegistration struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type SellerRegistration struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TCodeActivation struct {
	UserID         int          `json:"user_id"`
	CodeActivation string       `json:"code_activation"`
	IsUsed         string       `json:"is_used"`
	CreatedAt      time.Time    `json:"created_at"`
	UsedAt         sql.NullTime `json:"used_at"`
}
