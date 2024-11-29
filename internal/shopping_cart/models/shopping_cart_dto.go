package models

// request
type ShoppingCartReq struct {
	Id         string `json:"id" validate:"required,min=1,numeric"`
	CustomerId string `json:"customerId" validate:"required,min=1,numeric"`
	ProductId  string `json:"productId" validate:"required,min=1,numeric"`
	Quantity   string `json:"quantity" validate:"required,min=1,numeric"`
}
type DeleteReq struct {
	CustomerId string `json:"customerId" validate:"required,min=1,numeric"`
	ProductId  string `json:"productId" validate:"required,min=1,numeric"`
	CartId     string `json:"cartId" validate:"required,min=1,numeric"`
}

// response
type ShoppingCartRes struct {
	Id         string `json:"id,omitempty"`
	CustomerId string `json:"customerId,omitempty"`
	Status     string `json:"status,omitempty"`
}

type ShoppingCartItemRes struct {
	ProductId string `json:"productId,omitempty"`
	CartId    string `json:"cartId,omitempty"`
	Quantity  string `json:"quantity,omitempty"`
}

type ShoppingCartItem struct {
	ProductName string `json:"productName"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	NewCartId   string `json:"cartId,omitempty"`
}

type ListCartItemRes struct {
	Product    ShoppingCartItem `json:"product"`
	TotalPrice string           `json:"totalPrice"`
}
