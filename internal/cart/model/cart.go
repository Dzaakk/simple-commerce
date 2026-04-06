package model

import "time"

type Cart struct {
	ID         string
	CustomerID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
