package model

import "time"

type Order struct {
	ID              string
	OrderNumber     string
	CustomerID      string
	Status          string
	TotalAmount     float64
	ShippingAddress string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
