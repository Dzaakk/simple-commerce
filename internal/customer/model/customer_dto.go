package model

import (
	"Dzaakk/simple-commerce/package/template"
	"database/sql"
	"time"
)

type UpdateReq struct {
	CustomerID  string `json:"customer_id" validate:"required,min=1"`
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=02-01-2006"`
	Address     string `json:"address"`
}

func (req UpdateReq) ToCustomerModel(dateOfBirth time.Time, customerID int64) TCustomers {
	return TCustomers{
		ID:          customerID,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		DateOfBirth: sql.NullTime{Time: dateOfBirth, Valid: true},
		Address:     req.Address,
		Base: template.Base{
			UpdatedBy: sql.NullString{String: req.Username, Valid: true},
		},
	}
}

type CustomerRes struct {
	ID          string `json:"id"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Balance     string `json:"balance,omitempty"`
}
