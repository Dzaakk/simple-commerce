package model

import "time"

type Transaction struct {
	ID                string
	OrderID           string
	TransactionNumber string
	PaymentMethod     string
	Status            string
	Amount            float64
	PaidAt            *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
