package repositories

import (
	"Dzaakk/simple-commerce/internal/seller/models"
	"database/sql"
	"time"
)

type SellerRepositoryImpl struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) SellerRepository {
	return &SellerRepositoryImpl{
		DB: db,
	}
}

const (
	queryCreateSeller  = "INSERT INTO public.seller (name, email, password, balance, created, created_by) VALUES ($1, $2, $3, $4, $5, $6)"
	queryUpdateSeler   = "UPDATE public.seller SET name=$1, email=$2, password=$3, balance=$4, updated=NOW(), updated_by=$5 WHERE id=$6"
	queryFindById      = "SELECT * FROM public.seller WHERE id = $1"
	queryUpdateBalance = "UPDATE public.seller SET balance=$1, updated=NOW(), updated_by=$2 WHERE id=$3"
)

func (repo *SellerRepositoryImpl) Create(data models.SellerReq) error {
	statement, err := repo.DB.Prepare(queryCreateSeller)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(data.Name, data.Email, data.Password, 0, time.Now(), "SYSTEM")
	if err != nil {
		return err
	}

	return nil
}

// Update implements SellerRepository.
func (repo *SellerRepositoryImpl) Update(data models.SellerReq) error {
	panic("unimplemented")
}

// FindById implements SellerRepository.
func (repo *SellerRepositoryImpl) FindById(sellerId int64) (*models.TSeller, error) {
	panic("unimplemented")
}

// InsertBalance implements SellerRepository.
func (repo *SellerRepositoryImpl) InsertBalance(sellerId int64, balance int64) error {
	panic("unimplemented")
}
