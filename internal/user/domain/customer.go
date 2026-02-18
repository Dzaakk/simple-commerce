package domain

import "time"

type CustomerStatus string

const (
	CustomerStatusPending CustomerStatus = "pending"
	CustomerStatusActive  CustomerStatus = "active"
)

type Customer struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	Phone        string
	Status       CustomerStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
