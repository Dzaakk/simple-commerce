package model

import "time"

type CartItem struct {
	ID            int64
	CartID        string
	ProductID     string
	Quantity      int
	PriceSnapshot float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
