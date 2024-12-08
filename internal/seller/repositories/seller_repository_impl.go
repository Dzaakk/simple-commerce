package repositories

import (
	"Dzaakk/simple-commerce/internal/seller/models"
	model "Dzaakk/simple-commerce/internal/seller/models"
	"database/sql"
	"errors"
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
	queryUpdateSeler   = "UPDATE public.seller SET name=$1, email=$2, password=$3, updated=NOW(), updated_by=$4 WHERE id=$5"
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

func (repo *SellerRepositoryImpl) Update(data models.SellerReq) error {
	statement, err := repo.DB.Prepare(queryUpdateSeler)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(data.Name, data.Email, data.Password, time.Now(), data.Name)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SellerRepositoryImpl) FindById(sellerId int64) (*models.TSeller, error) {
	rows, err := repo.DB.Query(queryFindById, sellerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := retrieveData(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// InsertBalance implements SellerRepository.
func (repo *SellerRepositoryImpl) InsertBalance(sellerId int64, balance int64) error {
	panic("unimplemented")
}

func rowsToData(rows *sql.Rows) (*model.TSeller, error) {
	s := model.TSeller{}

	err := rows.Scan(&s.Id, &s.Name, &s.Email)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
func retrieveData(rows *sql.Rows) (*model.TSeller, error) {
	if rows.Next() {
		return rowsToData(rows)
	}
	return nil, errors.New("product not found")
}
