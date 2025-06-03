package model

// request
type ShoppingCartReq struct {
	ShoppingCartID string `json:"shopping_cart_id" validate:"required,min=1,numeric"`
	CustomerID     string `json:"customer_id" validate:"required,min=1,numeric"`
	ProductID      string `json:"product_id" validate:"required,min=1,numeric"`
	Quantity       string `json:"quantity" validate:"required,min=1,numeric"`
}
type DeleteReq struct {
	CustomerID string `json:"customer_id" validate:"required,min=1,numeric"`
	ProductID  string `json:"product_id" validate:"required,min=1,numeric"`
	CartID     string `json:"cart_id" validate:"required,min=1,numeric"`
}

// response
type ShoppingCartRes struct {
	ShoppingCartID string `json:"shooping_cart_id,omitempty"`
	CustomerID     string `json:"customer_id,omitempty"`
	Status         string `json:"status,omitempty"`
}

type ShoppingCartItemRes struct {
	ProductID string `json:"product_id,omitempty"`
	CartID    string `json:"cart_id,omitempty"`
	Quantity  string `json:"quantity,omitempty"`
}

type ShoppingCartItem struct {
	ProductName string `json:"product_name"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	NewCartID   string `json:"cart_id,omitempty"`
}

type ListCartItemRes struct {
	Product    ShoppingCartItem `json:"product"`
	TotalPrice string           `json:"totalPrice"`
}
