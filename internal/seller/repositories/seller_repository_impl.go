package repositories

import (
	model "Dzaakk/simple-commerce/internal/seller/models"
	"context"

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
	dbQueryTimeout      = 1 * time.Second
)

func (repo *SellerRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, dbQueryTimeout)
}

func (repo *SellerRepositoryImpl) Create(ctx context.Context, data model.TSeller) (int64, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	var id int64
	err := repo.DB.QueryRowContext(ctx, queryCreate, data.Username, data.Email, data.Password, 0, time.Now(), "SYSTEM").Scan(id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *SellerRepositoryImpl) Update(ctx context.Context, data model.TSeller) (int64, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()
	result, err := repo.DB.ExecContext(ctx, queryUpdate, data.Username, data.Email, time.Now(), data.Username)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *SellerRepositoryImpl) FindById(ctx context.Context, sellerId int64) (*model.TSeller, error) {
	if sellerId <= 0 {
		return nil, errors.New("invalid input parameter")
	}

	return repo.findSeller(ctx, queryFindById, sellerId)
}
func (repo *SellerRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.TSeller, error) {
	if username == "" {
		return nil, errors.New("invalid input parameter")
	}
	return repo.findSeller(ctx, queryFindByUsername, username)
}

func (repo *SellerRepositoryImpl) findSeller(ctx context.Context, query string, args ...interface{}) (*model.TSeller, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	row := repo.DB.QueryRowContext(ctx, query, args...)
	return scanSeller(row)
}

func (repo *SellerRepositoryImpl) InsertBalance(ctx context.Context, sellerId int64, balance int64) error {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryUpdateBalance, balance, sellerId)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SellerRepositoryImpl) UpdatePassword(ctx context.Context, sellerId int64, newPassword string) (int64, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryUpdatePassword, newPassword, sellerId)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *SellerRepositoryImpl) Deactive(ctx context.Context, sellerId int64) (int64, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryDeactive, "I", sellerId)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}
