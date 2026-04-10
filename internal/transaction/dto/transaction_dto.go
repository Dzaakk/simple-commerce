package dto

import (
	"Dzaakk/simple-commerce/package/constant"
	"time"
)

type CreateTransactionReq struct {
	CustomerID    string `json:"customer_id"`
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
}

type TransactionRes struct {
	ID                string                     `json:"id"`
	OrderID           string                     `json:"order_id"`
	TransactionNumber string                     `json:"transaction_number"`
	PaymentMethod     string                     `json:"payment_method"`
	Status            constant.TransactionStatus `json:"status"`
	Amount            float64                    `json:"amount"`
	PaidAt            *time.Time                 `json:"paid_at,omitempty"`
	CreatedAt         time.Time                  `json:"created_at"`
}

type PaymentCallbackReq struct {
	TransactionNumber string                    `json:"transaction_number"`
	Status            constant.TransactionStatus `json:"status"`
	PaidAt            *time.Time                `json:"paid_at,omitempty"`
	Signature         string                    `json:"signature"`
}
