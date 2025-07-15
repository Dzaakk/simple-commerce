package model

import (
	"Dzaakk/simple-commerce/package/template"
	"fmt"
)

type TProduct struct {
	ID          int     `json:"id"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  int     `json:"category_id"`
	SellerID    int     `json:"seller_id"`
	template.Base
}

func (p *TProduct) ToResponse() ProductRes {
	return ProductRes{
		ProductName: p.ProductName,
		Price:       fmt.Sprintf("%0.f", p.Price),
		Stock:       fmt.Sprintf("%d", p.Stock),
		CategoryID:  fmt.Sprintf("%d", p.CategoryID),
		SellerID:    fmt.Sprintf("%d", p.SellerID),
	}
}
