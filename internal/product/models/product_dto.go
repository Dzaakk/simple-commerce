package models

type ProductRes struct {
	ProductID   string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Price       string `json:"price,omitempty"`
	Stock       string `json:"stock,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
	SellerID    string `json:"seller_id,omitempty"`
}

type ProductReq struct {
	ProductID   string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Price       string `json:"price,omitempty"`
	Stock       string `json:"stock,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
	SellerID    string `json:"seller_id,omitempty"`
}
