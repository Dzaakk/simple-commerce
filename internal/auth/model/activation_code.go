package model

import "time"

type ActivationCode struct {
	ID        int64
	Code      string
	Email     string
	Type      string
	UserType  string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
