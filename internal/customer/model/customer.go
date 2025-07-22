package model

import (
	"Dzaakk/simple-commerce/package/template"
	"database/sql"
	"fmt"
)

type TCustomers struct {
	ID             int64        `json:"id"`
	Username       string       `json:"username"`
	Email          string       `json:"email"`
	PhoneNumber    string       `json:"phone_number"`
	Password       string       `json:"password"`
	Balance        float64      `json:"balance"`
	Status         uint8        `json:"status"`
	ProfilePicture string       `json:"profile_picture"`
	Address        string       `json:"address"`
	Gender         int          `json:"gender" validate:"number"`
	DateOfBirth    sql.NullTime `json:"date_of_birth"`
	LastLogin      sql.NullTime `json:"last_login"`
	template.Base
}

func (c *TCustomers) ToResponse() CustomerRes {
	return CustomerRes{
		ID:          fmt.Sprintf("%d", c.ID),
		Username:    c.Username,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
	}
}

type CustomerBalance struct {
	CustomerID int64   `json:"customer_id"`
	Balance    float64 `json:"balance"`
}
