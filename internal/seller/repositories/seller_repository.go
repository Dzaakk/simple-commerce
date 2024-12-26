package repositories

import model "Dzaakk/simple-commerce/internal/seller/models"

type SellerRepository interface {
	Create(data model.TSeller) (int64, error)
	FindById(sellerId int64) (*model.TSeller, error)
	FindByUsername(username string) (*model.TSeller, error)
	Update(data model.TSeller) (int64, error)
	InsertBalance(sellerId, balance int64) error
}
