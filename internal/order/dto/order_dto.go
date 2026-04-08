package dto

import (
	"Dzaakk/simple-commerce/package/constant"
	"time"
)

type CreateOrderReq struct {
	CustomerID      string         `json:"customer_id"`
	Items           []OrderItemReq `json:"items"`
	ShippingAddress string         `json:"shipping_address"`
}

type OrderItemReq struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderRes struct {
	ID              string               `json:"id"`
	OrderNumber     string               `json:"order_number"`
	Status          constant.OrderStatus `json:"status"`
	TotalAmount     float64              `json:"total_amount"`
	ShippingAddress string               `json:"shipping_address"`
	CreatedAt       time.Time            `json:"created_at"`
}

type OrderItemRes struct {
	ProductID string  `json:"product_id"`
	SellerID  string  `json:"seller_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

type OrderDetailRes struct {
	OrderRes
	Items []OrderItemRes `json:"items"`
}
