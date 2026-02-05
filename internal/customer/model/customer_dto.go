package model

import (
	"Dzaakk/simple-commerce/package/template"
	"database/sql"
	"time"
)

type CreateReq struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=02-01-2006"`
	Address     string `json:"address"`
}

func (c *CreateReq) ToCreateData(dateOfBirth time.Time) *TCustomers {
	return &TCustomers{
		Username:    c.Username,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
		Password:    c.Password,
		Balance:     float64(10000000),
		Status:      1,
		// Gender:         gender,
		DateOfBirth:    sql.NullTime{Valid: true, Time: dateOfBirth},
		LastLogin:      sql.NullTime{Time: time.Now(), Valid: true},
		ProfilePicture: "",
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}
}

type UpdateReq struct {
	CustomerID  string `json:"customer_id" validate:"required,min=1"`
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=02-01-2006"`
	Address     string `json:"address"`
}

func (u *UpdateReq) ToUpdateData(dateOfBirth time.Time, customerID int64) *TCustomers {
	return &TCustomers{
		ID:          customerID,
		Username:    u.Username,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		DateOfBirth: sql.NullTime{Time: dateOfBirth, Valid: true},
		Address:     u.Address,
		Base: template.Base{
			UpdatedBy: sql.NullString{String: u.Username, Valid: true},
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
