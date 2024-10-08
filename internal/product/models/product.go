package product

import "Dzaakk/simple-commerce/package/template"

type TProduct struct {
	Id          int     `json:"int"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryId  int     `json:"category_id"`
	template.Base
}

type ProductRes struct {
	Id          string `json:"int,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Price       string `json:"price,omitempty"`
	Stock       string `json:"stock,omitempty"`
	CategoryId  string `json:"category_id,omitempty"`
}
