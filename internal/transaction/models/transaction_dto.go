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

type CustomerListTransactionRes struct {
	CustomerId      string                `json:"customerId"`
	ListTransaction []CustomerTransaction `json:"listTransaction"`
}
