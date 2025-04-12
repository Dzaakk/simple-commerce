package models

import template "Dzaakk/simple-commerce/package/templates"

type THistoryTransaction struct {
	TransactionID int64   `json:"transaction_id"`
	CustomerID    int64   `json:"customer_id"`
	ProductName   string  `json:"productName"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	Status        string  `json:"status"`
	template.Base
}
type HistoryTransaction struct {
	TransactionID string `json:"transaction_id"`
	CustomerID    string `json:"customer_id"`
	ProductName   string `json:"productName"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	Status        string `json:"status"`
}
