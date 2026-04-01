package model

import "time"

type Inventory struct {
	ID               int64
	ProductID        string
	StockQuantity    int
	ReservedQuantity int
	Version          int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
