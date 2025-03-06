package repositories

import (
	model "Dzaakk/simple-commerce/internal/auth/models"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
	"time"
)

type AuthRepositoryImpl struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &AuthRepositoryImpl{DB: db}
}

const (
	queryCreate           = `INSERT INTO public.code_activation (customer_id, code_activation, is_used, created_at, used_at) VALUES ($1, $2, $3, $4, $5)`
	queryFindByCustomerID = `SELECT * FROM public.code_activation WHERE customer_id = $1`
	dbQueryTimeout        = 3 * time.Second
)

func (repo *AuthRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, dbQueryTimeout)
}

func (repo *AuthRepositoryImpl) CreateCode(c context.Context, data model.TActivationCode) error {
	ctx, cancel := repo.contextWithTimeout(c)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryCreate, data.CustomerID, data.CodeActivation, data.IsUsed, data.CreatedAt, data.UsedAt)
	if err != nil {
		return response.ExecError("create activation code", err)
	}

	return nil
}

func (repo *AuthRepositoryImpl) FindCodeByCustomerId(c context.Context, id int64) (*model.TActivationCode, error) {
	ctx, cancel := repo.contextWithTimeout(c)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, queryFindByCustomerID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activationCode, err := retrieveCodeActivaton(rows)
	if err != nil {
		return nil, err
	}

	return activationCode, nil
}

func rowsToActivationCode(rows *sql.Rows) (*model.TActivationCode, error) {
	ac := model.TActivationCode{}
	err := rows.Scan(&ac.CustomerID, &ac.CodeActivation, &ac.IsUsed, &ac.CreatedAt, &ac.UsedAt)
	if err != nil {
		return nil, err
	}

	return &ac, nil
}

func retrieveCodeActivaton(rows *sql.Rows) (*model.TActivationCode, error) {
	if rows.Next() {
		return rowsToActivationCode(rows)
	}
	return nil, errors.New("code activation not found")
}
