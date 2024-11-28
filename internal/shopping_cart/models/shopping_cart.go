package models

import template "Dzaakk/simple-commerce/package/templates"

type TShoppingCart struct {
	Id         int    `json:"id"`
	CustomerId int    `json:"customer_id"`
	Status     string `json:"status"`
	template.Base
}
type TShoppingCartItem struct {
	ProductId int `json:"productId"`
	CartId    int `json:"cart_id"`
	Quantity  int `json:"quantity"`
	template.Base
}
type TCartItemDetail struct {
	ProductId   int     `json:"productId"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
}
