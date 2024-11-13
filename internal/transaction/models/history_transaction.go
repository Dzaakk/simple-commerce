package models

import "Dzaakk/simple-commerce/package/template"

type THistoryTransaction struct {
	Id          int64   `json:"id"`
	CustomerId  int64   `json:"customer_id"`
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Status      string  `json:"status"`
	template.Base
}
type HistoryTransaction struct {
	Id          string `json:"id"`
	CustomerId  string `json:"customer_id"`
	ProductName string `json:"productName"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	Status      string `json:"status"`
}
