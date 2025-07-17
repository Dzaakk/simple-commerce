package model

import (
	"Dzaakk/simple-commerce/package/template"
	"time"
)

type TOrder struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	OrderData  string    `json:"order_data"`
	ExpiredAt  time.Time `json:"expired_at"`
	Paid       bool      `json:"paid"`
	Base       template.Base
}
