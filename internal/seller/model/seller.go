package model

import "Dzaakk/simple-commerce/package/template"

type TSeller struct {
	Id       int64   `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
	template.Base
}
