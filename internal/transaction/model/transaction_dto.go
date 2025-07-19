package model

type TransactionReq struct {
	CustomerID string `json:"customerId" validate:"required,min=1,numeric"`
	CartID     string `json:"cartId" validate:"required,min=1,numeric"`
	OrderID    string `json:"orderId" validate:"required,min=1,numeric"`
}
type TransactionRes struct {
	CustomerID       string        `json:"customerId"`
	TransactionDate  string        `json:"transactionDate"`
	ListProduct      []ListProduct `json:"listProduct"`
	TotalTransaction string        `json:"totalTransaction"`
}

type CustomerListTransactionRes struct {
	CustomerID      string                `json:"customerId"`
	ListTransaction []CustomerTransaction `json:"listTransaction"`
}
