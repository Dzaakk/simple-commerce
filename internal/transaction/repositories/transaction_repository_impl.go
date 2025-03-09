package repositories

import (
	model "Dzaakk/simple-commerce/internal/transaction/models"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"time"
)

type TransactionRepositoryImpl struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &TransactionRepositoryImpl{
		DB: db,
	}
}

const (
	queryCreateTransaction = `INSERT INTO public.transaction (customer_id, cart_id, total_amount, transaction_date, status, created, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	dbQueryTimeout         = 3 * time.Second
)

func (repo *TransactionRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, dbQueryTimeout)
}

func (repo *TransactionRepositoryImpl) Create(ctx context.Context, data model.TTransaction) (*model.TTransaction, error) {
	c, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(c, queryCreateTransaction, data.CustomerId, data.CartId, data.TotalAmount, data.TransactionDate, data.Status, data.Base.Created, data.Base.CreatedBy)
	if err != nil {
		return nil, response.ExecError("create transaction", err)
	}

	id, _ := result.LastInsertId()
	data.Id = int(id)

	return &data, nil
}
func (repo *TransactionRepositoryImpl) CreateWithTx(ctx context.Context, tx *sql.Tx, data model.TTransaction) (*model.TTransaction, error) {
	c, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	statement, err := tx.Prepare(queryCreateTransaction)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var id int

	err = statement.QueryRowContext(c, data.CustomerId, data.CartId, data.TotalAmount, data.TransactionDate, data.Status, data.Base.Created, data.Base.CreatedBy).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *TransactionRepositoryImpl) BeginTransaction() (*sql.Tx, error) {
	return repo.DB.Begin()
}
