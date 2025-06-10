package model

import (
	"database/sql"
)

type CustomerToken struct {
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	PhoneNumber string       `json:"phone_number"`
	Password    string       `json:"password"`
	Gender      int          `json:"gender"`
	DateOfBirth sql.NullTime `json:"date_of_birth"`
}

type SellerToken struct {
}
