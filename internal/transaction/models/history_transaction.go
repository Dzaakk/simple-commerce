package models

import template "Dzaakk/simple-commerce/package/templates"

type THistoryTransaction struct {
	TransactionId int64   `json:"transaction_id"`
	CustomerId    int64   `json:"customer_id"`
	ProductName   string  `json:"productName"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	Status        string  `json:"status"`
	template.Base
}
type HistoryTransaction struct {
	TransactionId string `json:"transaction_id"`
	CustomerId    string `json:"customer_id"`
	ProductName   string `json:"productName"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	Status        string `json:"status"`
}
