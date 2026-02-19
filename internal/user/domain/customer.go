package domain

import (
	"time"
)

type Customer struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	Phone        string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
