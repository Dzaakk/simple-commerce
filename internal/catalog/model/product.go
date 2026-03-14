package model

import "time"

type Product struct {
	ID          string
	SellerID    string
	CategoryID  int64
	Name        string
	SKU         string
	Description *string
	Price       float64
	ImageURL    *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
