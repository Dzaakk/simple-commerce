package model

import "time"

type OrderItem struct {
	ID        int64
	OrderID   string
	ProductID string
	SellerID  string
	Quantity  int
	Price     float64
	Subtotal  float64
	CreatedAt time.Time
}
