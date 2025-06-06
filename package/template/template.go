package template

import (
	"database/sql"
	"time"
)

type Base struct {
	Created   time.Time      `json:"created,omitempty"`
	CreatedBy string         `json:"createdBy,omitempty"`
	Updated   sql.NullTime   `json:"updated,omitempty"`
	UpdatedBy sql.NullString `json:"updatedBy,omitempty"`
}
