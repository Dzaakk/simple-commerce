package dto

import (
	"Dzaakk/simple-commerce/internal/user/model"
	"fmt"
	"time"
)

type CreateReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone"`
}

func (c *CreateReq) ToCreateData() *model.Customer {
	return &model.Customer{
		Email:        c.Email,
		PasswordHash: c.Password,
		FullName:     c.FullName,
		Phone:        c.Phone,
		Status:       "pending",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

type UpdateReq struct {
	CustomerID string `json:"customer_id" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	FullName   string `json:"full_name" validate:"required"`
	Phone      string `json:"phone"`
	Status     string `json:"status"`
}

func (u *UpdateReq) ToUpdateData(customerID int64) *model.Customer {
	status := u.Status
	if status == "" {
		status = "pending"
	}

	return &model.Customer{
		ID:        fmt.Sprintf("%d", customerID),
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Status:    status,
		UpdatedAt: time.Now(),
	}
}

type CustomerRes struct {
	ID        string    `json:"id"`
	Email     string    `json:"email,omitempty"`
	FullName  string    `json:"full_name,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func ToCustomerRes(c *model.Customer) CustomerRes {
	return CustomerRes{
		ID:        c.ID,
		Email:     c.Email,
		FullName:  c.FullName,
		Phone:     c.Phone,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
