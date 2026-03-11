package model

import "time"

type RefreshToken struct {
	ID        int64
	UserID    string
	UserType  string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
