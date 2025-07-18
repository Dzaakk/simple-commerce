package template

import (
	"database/sql"
	"time"
)

const (
	StatusActive   = "A"
	StatusInactive = "I"
	StatusBlock    = "B"
	RoleCustomer   = "CUSTOMER"
	RoleSeller     = "SELLER"
	FormatDate     = "02-01-2006"
)

type Base struct {
	Created   time.Time      `json:"created,omitempty"`
	CreatedBy string         `json:"createdBy,omitempty"`
	Updated   sql.NullTime   `json:"updated,omitempty"`
	UpdatedBy sql.NullString `json:"updatedBy,omitempty"`
}
