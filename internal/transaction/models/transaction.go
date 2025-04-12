package models

import (
	template "Dzaakk/simple-commerce/package/templates"
	"time"
)

type TTransaction struct {
	ID              int       `json:"id"`
	CustomerID      int       `json:"customer_id"`
	CartID          int       `json:"cart_id"`
	TotalAmount     float32   `json:"total_amount"`
	TransactionDate time.Time `json:"transaction_date"`
	Status          string    `json:"status"`
	template.Base
}

type TCartItemDetail struct {
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type ListProduct struct {
	ProductName string `json:"product_name"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
}

type CustomerTransaction struct {
	ProductName     string `json:"product_name"`
	Price           string `json:"price"`
	Quantity        string `json:"quantity"`
	TransactionDate string `json:"transactionDate"`
}
