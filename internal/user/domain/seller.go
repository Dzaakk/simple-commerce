package domain

import "time"

type Seller struct {
	ID           string
	Email        string
	PasswordHash string
	ShopName     string
	Phone        string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
