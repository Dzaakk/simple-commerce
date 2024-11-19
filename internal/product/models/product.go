package models

import template "Dzaakk/simple-commerce/package/templates"

type TProduct struct {
	Id          int     `json:"id,omitempty,string"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price,string"`
	Stock       int     `json:"stock,string"`
	CategoryId  int     `json:"category_id,string"`
	SellerId    int     `json:"seller_id,string"`
	template.Base
}

type ProductRes struct {
	Id          string `json:"id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Price       string `json:"price,omitempty"`
	Stock       string `json:"stock,omitempty"`
	CategoryId  string `json:"category_id,omitempty"`
}
