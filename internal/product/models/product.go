package models

import template "Dzaakk/simple-commerce/package/templates"

type TProduct struct {
	Id          int     `json:"id,omitempty"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryId  int     `json:"category_id"`
	SellerId    int     `json:"seller_id"`
	template.Base
}
