package models

import "Dzaakk/simple-commerce/package/templates"

type TSeller struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Balance  float64 `json:"balance"`
	templates.Base
}

type SellerReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SellerRes struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
