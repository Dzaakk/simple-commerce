package transaction

import (
	model "Dzaakk/simple-commerce/internal/transaction/models"
	"database/sql"
)

type TransactionRepositoryImpl struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &TransactionRepositoryImpl{
		DB: db,
	}
}

const queryCreateTransaction = `INSERT INTO public.transaction (customer_id, cart_id, total_amount, transaction_date, status, created, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

func (repo *TransactionRepositoryImpl) Create(data model.TTransaction) (*model.TTransaction, error) {
	statement, err := repo.DB.Prepare(queryCreateTransaction)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var id int

	err = statement.QueryRow(data.CustomerId, data.CartId, data.TotalAmount, data.TransactionDate, data.Status, data.Base.Created, data.Base.CreatedBy).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
