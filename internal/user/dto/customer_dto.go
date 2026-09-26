package dto

import (
	"time"

	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
)

type RegisterCustomerRequest struct {
	Email        string
	PasswordHash string
	FullName     string
}

func (r RegisterCustomerRequest) ToModel() *model.Customer {
	now := time.Now()
	return &model.Customer{
		Email: r.Email, PasswordHash: r.PasswordHash, FullName: r.FullName,
		Status: string(constant.StatusActive), CreatedAt: now, UpdatedAt: now,
	}
}

type CustomerResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCustomerResponse(customer *model.Customer) CustomerResponse {
	return CustomerResponse{
		ID: customer.ID, Email: customer.Email, FullName: customer.FullName,
		Status: customer.Status, CreatedAt: customer.CreatedAt, UpdatedAt: customer.UpdatedAt,
	}
}
