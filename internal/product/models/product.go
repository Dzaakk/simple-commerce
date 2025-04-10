package models

import template "Dzaakk/simple-commerce/package/templates"

type TProduct struct {
	ProductID   int     `json:"id"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  int     `json:"category_id"`
	SellerID    int     `json:"seller_id"`
	template.Base
}
