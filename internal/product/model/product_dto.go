package model

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

type ProductFilter struct {
	ProductName string `form:"product_name"`
	CategoryID  string `form:"category_id"`
	SellerID    string `form:"seller_id"`
	LowPrice    string `form:"low_price"`
	HighPrice   string `form:"high_price"`
	Offset      int    `form:"offset"`
	Limit       int    `form:"limit"`
}
