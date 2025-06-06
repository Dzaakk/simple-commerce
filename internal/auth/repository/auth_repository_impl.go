package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"time"
)

type AuthRepositoryImpl struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &AuthRepositoryImpl{DB: db}
}

const (
	queryCreateCustomerCode = `INSERT INTO public.customer_activation_code (customer_id, code_activation, is_used, created_at) VALUES ($1, $2, $3, $4)`
	queryFindByCustomerID   = `SELECT * FROM public.customer_activation_code WHERE customer_id = $1`
	queryCreateSellerCode   = `INSERT INTO public.seller_activation_code (seller_id, code_activation, is_used, created_at) VALUES ($1, $2, $3, $4)`
	queryFindBySellerID     = `SELECT * FROM public.seller_activation_code WHERE seller_id = $1`
	dbQueryTimeout          = 3 * time.Second
)

func (repo *AuthRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, dbQueryTimeout)
}

func (repo *AuthRepositoryImpl) InsertCustomerCodeActivation(c context.Context, data model.TCustomerActivationCode) error {
	ctx, cancel := repo.contextWithTimeout(c)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryCreateCustomerCode, data.CustomerID, data.CodeActivation, data.IsUsed, data.CreatedAt)
	if err != nil {
		return response.ExecError("create activation code", err)
	}

	return nil
}

func (repo *AuthRepositoryImpl) FindCodeByCustomerID(c context.Context, customerID int64) (*model.TCustomerActivationCode, error) {
	ctx, cancel := repo.contextWithTimeout(c)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, queryFindByCustomerID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activationCode, err := retrieveCustomerCodeActivaton(rows)
	if err != nil {
		return nil, err
	}

	return activationCode, nil
}

func (repo *AuthRepositoryImpl) FindCodeBySellerID(c context.Context, sellerID int64) (*model.TSellerActivationCode, error) {
	ctx, cancel := repo.contextWithTimeout(c)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, queryFindBySellerID, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activationCode, err := retrieveSellerCodeActivaton(rows)
	if err != nil {
		return nil, err
	}

	return activationCode, nil
}

func (repo *AuthRepositoryImpl) InsertSellerCodeActivation(c context.Context, data model.TSellerActivationCode) error {
	ctx, cancel := repo.contextWithTimeout(c)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryCreateSellerCode, data.SellerID, data.CodeActivation, data.IsUsed, data.CreatedAt)
	if err != nil {
		return response.ExecError("create activation code", err)
	}

	return nil
}
