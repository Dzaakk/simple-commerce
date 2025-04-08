package models

import template "Dzaakk/simple-commerce/package/templates"

type TCustomers struct {
	ID          int     `json:"customer_id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	Password    string  `json:"password"`
	Balance     float64 `json:"balance"`
	Status      string  `json:"status"`
	template.Base
}
type CustomerBalance struct {
	ID      int64   `json:"customer_id"`
	Balance float64 `json:"balance"`
}
