package models

import template "Dzaakk/simple-commerce/package/templates"

type TCustomers struct {
	Id          int     `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	Password    string  `json:"password"`
	Balance     float64 `json:"balance"`
	Status      string  `json:"status"`
	template.Base
}
type CustomerBalance struct {
	Id      int     `json:"id"`
	Balance float64 `json:"balance"`
}
