package model

import "time"

type Category struct {
	ID        int64
	Name      string
	Slug      string
	CreatedAt time.Time
}
