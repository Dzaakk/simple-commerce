package models

import "Dzaakk/simple-commerce/package/template"

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

// request
type ShoppingCartReq struct {
	Id         string `json:"id,omitempty"`
	CustomerId string `json:"customerId,omitempty"`
	ProductId  string `json:"productId,omitempty"`
	Quantity   string `json:"quantity,omitempty"`
}
type DeleteReq struct {
	CustomerId string `json:"customerId,omitempty"`
	CartId     string `json:"cartId,omitempty"`
	ProductId  string `json:"productId,omitempty"`
}

// response
type ShoppingCartRes struct {
	Id         string `json:"id,omitempty"`
	CustomerId string `json:"customerId,omitempty"`
	Status     string `json:"status,omitempty"`
}

type ShoppingCartItemRes struct {
	ProductId string `json:"productId,omitempty"`
	CartId    string `json:"cartId,omitempty"`
	Quantity  string `json:"quantity,omitempty"`
}

type ShoppingCartItem struct {
	ProductName string `json:"productName"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	NewCartId   string `json:"cartId,omitempty"`
}

type ListCartItemRes struct {
	Product    ShoppingCartItem `json:"product"`
	TotalPrice string           `json:"totalPrice"`
}
