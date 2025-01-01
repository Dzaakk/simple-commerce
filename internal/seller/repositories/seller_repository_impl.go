package repositories

import (
	model "Dzaakk/simple-commerce/internal/seller/models"
	// template "Dzaakk/simple-commerce/package/templates"
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
	queryCreate         = "INSERT INTO public.seller (username, email, password, balance, created, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	queryUpdate         = "UPDATE public.seller SET username=$1, email=$2, updated=NOW(), updated_by=$3 WHERE id=$4"
	queryUpdatePassword = "UPDATE public.seller set password=$1 WHERE id=$2"
	queryDeactive       = "UPDATE public.seller set status=$1 WHERE id=$2"
	queryFindById       = "SELECT * FROM public.seller WHERE id = $1"
	queryFindByUsername = "SELECT * FROM public.seller WHERE username = $1"
	queryUpdateBalance  = "UPDATE public.seller SET balance=$1, updated=NOW(), updated_by=$2 WHERE id=$2"
)

func (repo *SellerRepositoryImpl) Create(data model.TSeller) (int64, error) {
	statement, err := repo.DB.Prepare(queryCreate)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var id int64
	err = statement.QueryRow(data.Username, data.Email, data.Password, 0, time.Now(), "SYSTEM").Scan(id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *SellerRepositoryImpl) Update(data model.TSeller) (int64, error) {
	statement, err := repo.DB.Prepare(queryUpdate)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(data.Username, data.Email, time.Now(), data.Username)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *SellerRepositoryImpl) FindById(sellerId int64) (*model.TSeller, error) {
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

func (repo *SellerRepositoryImpl) InsertBalance(sellerId int64, balance int64) error {
	statement, err := repo.DB.Prepare(queryUpdateBalance)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(balance, sellerId)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SellerRepositoryImpl) FindByUsername(username string) (*model.TSeller, error) {
	rows, err := repo.DB.Query(queryFindByUsername, username)
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

func (repo *SellerRepositoryImpl) UpdatePassword(sellerId int64, newPassword string) (int64, error) {
	statement, err := repo.DB.Prepare(queryUpdatePassword)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(newPassword, sellerId)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *SellerRepositoryImpl) Deactive(sellerId int64) (int64, error) {
	statement, err := repo.DB.Prepare(queryUpdatePassword)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec("I", sellerId)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func rowsToData(rows *sql.Rows) (*model.TSeller, error) {
	s := model.TSeller{}
	// b := template.Base{}
	err := rows.Scan(&s.Id, &s.Username, &s.Email, &s.Balance, &s.Password, &s.Created, s.CreatedBy, s.Updated, s.UpdatedBy)
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
