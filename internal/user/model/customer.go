package model

import "time"

type Customer struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
