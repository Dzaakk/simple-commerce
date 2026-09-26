package model

import "time"

type Product struct {
	ID          int64
	CategoryID  int64
	SKU         string
	Name        string
	Description string
	Price       float64
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
