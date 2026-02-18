package domain

import "time"

type SellerStatus string

const (
	SellerStatusPending SellerStatus = "pending"
	SellerStatusActive  SellerStatus = "active"
)

type Seller struct {
	ID           string
	Email        string
	PasswordHash string
	ShopName     string
	Phone        string
	Status       SellerStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
