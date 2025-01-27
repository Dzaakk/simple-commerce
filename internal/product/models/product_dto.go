package models

type ProductRes struct {
	Id          string `json:"id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Price       string `json:"price,omitempty"`
	Stock       string `json:"stock,omitempty"`
	CategoryId  string `json:"category_id,omitempty"`
	SellerId    string `json:"seller_id,omitempty"`
}

type ProductReq struct {
	Id          string `json:"id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Price       string `json:"price,omitempty"`
	Stock       string `json:"stock,omitempty"`
	CategoryId  string `json:"category_id,omitempty"`
	SellerId    string `json:"seller_id,omitempty"`
}
