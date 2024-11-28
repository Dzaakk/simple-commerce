package models

type TransactionReq struct {
	CustomerId string `json:"customerId" validate:"required,min=1,numeric"`
	CartId     string `json:"cartId" validate:"required,min=1,numeric"`
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
