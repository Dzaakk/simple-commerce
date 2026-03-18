package model

import (
	"Dzaakk/simple-commerce/package/constant"
	"time"
)

type RefreshToken struct {
	ID        int64
	UserID    string
	UserType  constant.UserType
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
