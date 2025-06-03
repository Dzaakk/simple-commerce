package model

import "Dzaakk/simple-commerce/package/template"

type TProduct struct {
	ID          int     `json:"id"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  int     `json:"category_id"`
	SellerID    int     `json:"seller_id"`
	template.Base
}
