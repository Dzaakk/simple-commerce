package dto

import (
	"Dzaakk/simple-commerce/internal/user/domain"
	"fmt"
	"time"
)

type SellerCreateReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	ShopName string `json:"shop_name" validate:"required"`
	Phone    string `json:"phone"`
}

func (s *SellerCreateReq) ToCreateData() *domain.Seller {
	return &domain.Seller{
		Email:        s.Email,
		PasswordHash: s.Password,
		ShopName:     s.ShopName,
		Phone:        s.Phone,
		Status:       "pending",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

type SellerUpdateReq struct {
	SellerID string `json:"seller_id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	ShopName string `json:"shop_name" validate:"required"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

func (s *SellerUpdateReq) ToUpdateData(sellerID int64) *domain.Seller {
	status := s.Status
	if status == "" {
		status = "pending"
	}

	return &domain.Seller{
		ID:        fmt.Sprintf("%d", sellerID),
		Email:     s.Email,
		ShopName:  s.ShopName,
		Phone:     s.Phone,
		Status:    status,
		UpdatedAt: time.Now(),
	}
}

type SellerRes struct {
	ID        string    `json:"id"`
	Email     string    `json:"email,omitempty"`
	ShopName  string    `json:"shop_name,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func ToSellerRes(s *domain.Seller) SellerRes {
	return SellerRes{
		ID:        s.ID,
		Email:     s.Email,
		ShopName:  s.ShopName,
		Phone:     s.Phone,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
