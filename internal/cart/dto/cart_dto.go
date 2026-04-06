package dto

type CartItemRes struct {
	ProductID     string  `json:"product_id"`
	Quantity      int     `json:"quantity"`
	PriceSnapshot float64 `json:"price_snapshot"`
	Subtotal      float64 `json:"subtotal"`
}

type CartRes struct {
	CartID     string        `json:"cart_id"`
	CustomerID string        `json:"customer_id"`
	Items      []CartItemRes `json:"items"`
	Total      float64       `json:"total"`
}
