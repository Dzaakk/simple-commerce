package models

// request
type ShoppingCartReq struct {
	Id         string `json:"id,omitempty"`
	CustomerId string `json:"customerId,omitempty"`
	ProductId  string `json:"productId,omitempty"`
	Quantity   string `json:"quantity,omitempty"`
}
type DeleteReq struct {
	CustomerId string `json:"customerId,omitempty"`
	CartId     string `json:"cartId,omitempty"`
	ProductId  string `json:"productId,omitempty"`
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
