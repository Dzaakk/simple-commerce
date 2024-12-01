package repositories

import (
	"Dzaakk/simple-commerce/internal/seller/models"
	"database/sql"
)

type SellerRepositoryImpl struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) SellerRepository {
	return &SellerRepositoryImpl{
		DB: db,
	}
}

// Create implements SellerRepository.
func (s *SellerRepositoryImpl) Create(models.SellerReq) error {
	panic("unimplemented")
}

// FindById implements SellerRepository.
func (s *SellerRepositoryImpl) FindById(sellerId int64) (*models.TSeller, error) {
	panic("unimplemented")
}

// InsertBalance implements SellerRepository.
func (s *SellerRepositoryImpl) InsertBalance(sellerId int64, balance int64) error {
	panic("unimplemented")
}

// Update implements SellerRepository.
func (s *SellerRepositoryImpl) Update(models.SellerReq) error {
	panic("unimplemented")
}
