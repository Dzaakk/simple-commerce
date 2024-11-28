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
