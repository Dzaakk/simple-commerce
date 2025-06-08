package model

import (
	"Dzaakk/simple-commerce/package/template"
	"database/sql"
)

type TCustomers struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	PhoneNumber    string  `json:"phone_number"`
	Password       string  `json:"password"`
	Balance        float64 `json:"balance"`
	Status         int     `json:"status"`
	ProfilePicture string
	Gender         int          `json:"gender" validate:"number"`
	DateOfBirth    sql.NullTime `json:"date_of_birth"`
	LastLogin      sql.NullTime `json:"last_login"`
	template.Base
}
type CustomerBalance struct {
	CustomerID int64   `json:"customer_id"`
	Balance    float64 `json:"balance"`
}
