package models

import "Dzaakk/simple-commerce/package/template"

type TTransactionItem struct {
	Id          int64  `json:"id"`
	CustomerId  int64  `json:"customer_id"`
	ProductName string `json:"productName"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	Status      string `json:"status"`
	template.Base
}
