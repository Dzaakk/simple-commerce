package template

import (
	"database/sql"
	"time"
)

type Base struct {
	Created   time.Time      `json:"created"`
	CreatedBy string         `json:"createdBy"`
	Updated   sql.NullTime   `json:"updated,omitempty"`
	UpdatedBy sql.NullString `json:"updatedBy,omitempty"`
}
