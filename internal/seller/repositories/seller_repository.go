package repositories

import model "Dzaakk/simple-commerce/internal/seller/models"

type SellerRepository interface {
	Create(model.SellerReq) error
	FindById(sellerId int64) (*model.TSeller, error)
	Update(model.SellerReq) error
	InsertBalance(sellerId, balance int64) error
}
