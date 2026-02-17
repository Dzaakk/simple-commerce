package model

import (
	"time"
)

type Customers struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (c *Customers) ToResponse() CustomerRes {
	return CustomerRes{
		ID:       c.ID,
		FullName: c.FullName,
		Email:    c.Email,
		Phone:    c.Phone,
	}
}

type CustomerBalance struct {
	CustomerID string  `json:"customer_id"`
	Balance    float64 `json:"balance"`
}
