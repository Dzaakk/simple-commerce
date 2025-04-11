package models

import template "Dzaakk/simple-commerce/package/templates"

type TShoppingCart struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	Status     string `json:"status"`
	template.Base
}
type TShoppingCartItem struct {
	ProductID int `json:"product_id"`
	CartID    int `json:"cart_id"`
	Quantity  int `json:"quantity"`
	template.Base
}
type TCartItemDetail struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
}
