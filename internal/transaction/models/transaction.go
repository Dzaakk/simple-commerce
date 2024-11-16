package models

import (
	template "Dzaakk/simple-commerce/package/templates"
	"time"
)

type TTransaction struct {
	Id              int       `json:"id"`
	CustomerId      int       `json:"customer_id"`
	CartId          int       `json:"cart_id"`
	TotalAmount     float32   `json:"total_amount"`
	TransactionDate time.Time `json:"transaction_date"`
	Status          string    `json:"status"`
	template.Base
}

type TransactionReq struct {
	CustomerId string `json:"customerId"`
	CartId     string `json:"cartId"`
}
type TransactionRes struct {
	CustomerId       string        `json:"customerId"`
	TransactionDate  string        `json:"transactionDate"`
	ListProduct      []ListProduct `json:"listProduct"`
	TotalTransaction string        `json:"totalTransaction"`
}
type ListProduct struct {
	ProductName string `json:"productName"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
}

type TCartItemDetail struct {
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type CustomerListTransactionRes struct {
	CustomerId      string                `json:"customerId"`
	ListTransaction []CustomerTransaction `json:"listTransaction"`
}

type CustomerTransaction struct {
	ProductName     string `json:"product_name"`
	Price           string `json:"price"`
	Quantity        string `json:"quantity"`
	TransactionDate string `json:"transactionDate"`
}
