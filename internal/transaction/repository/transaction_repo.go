package repository

import (
	"Dzaakk/simple-commerce/internal/transaction/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	transactionSelectColumns = "id, order_id, transaction_number, payment_method, status, amount, paid_at, created_at, updated_at"

	transactionQueryCreate                  = "INSERT INTO public.transactions (order_id, transaction_number, payment_method, status, amount, paid_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	transactionQueryFindByID                = "SELECT " + transactionSelectColumns + " FROM public.transactions WHERE id=$1"
	transactionQueryFindByOrderID           = "SELECT " + transactionSelectColumns + " FROM public.transactions WHERE order_id=$1"
	transactionQueryFindByTransactionNumber = "SELECT " + transactionSelectColumns + " FROM public.transactions WHERE transaction_number=$1"
	transactionQueryUpdateStatus            = "UPDATE public.transactions SET status=$1, paid_at=$2, updated_at=$3 WHERE id=$4"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *sql.Tx, data *model.Transaction) (string, error) {
	if tx == nil {
		return "", errors.New("transaction is required")
	}

	var id string
	err := tx.QueryRowContext(
		ctx,
		transactionQueryCreate,
		data.OrderID,
		data.TransactionNumber,
		data.PaymentMethod,
		data.Status,
		data.Amount,
		data.PaidAt,
		data.CreatedAt,
		data.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return "", response.Error("failed to create transaction", err)
	}

	return id, nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, transactionID string) (*model.Transaction, error) {
	row := r.DB.QueryRowContext(ctx, transactionQueryFindByID, transactionID)

	return scanTransaction(row)
}

func (r *TransactionRepository) FindByOrderID(ctx context.Context, orderID string) (*model.Transaction, error) {
	row := r.DB.QueryRowContext(ctx, transactionQueryFindByOrderID, orderID)

	return scanTransaction(row)
}

func (r *TransactionRepository) FindByTransactionNumber(ctx context.Context, txNumber string) (*model.Transaction, error) {
	row := r.DB.QueryRowContext(ctx, transactionQueryFindByTransactionNumber, txNumber)

	return scanTransaction(row)
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, transactionID string, status constant.TransactionStatus, paidAt *time.Time) error {
	if tx == nil {
		return errors.New("transaction is required")
	}

	result, err := tx.ExecContext(ctx, transactionQueryUpdateStatus, status, paidAt, time.Now(), transactionID)
	if err != nil {
		return response.ExecError("update transaction status", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return response.Error("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return response.Error("no rows updated", sql.ErrNoRows)
	}

	return nil
}

func (r *TransactionRepository) GenerateTransactionNumber(ctx context.Context) (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	var count int
	err := r.DB.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM public.transactions WHERE created_at >= $1 AND created_at < $2",
		start,
		end,
	).Scan(&count)
	if err != nil {
		return "", response.Error("failed to count transactions", err)
	}

	seq := count + 1
	return fmt.Sprintf("TRX-%s-%04d", dateStr, seq), nil
}

func scanTransaction(row *sql.Row) (*model.Transaction, error) {
	var tx model.Transaction

	if err := row.Scan(
		&tx.ID,
		&tx.OrderID,
		&tx.TransactionNumber,
		&tx.PaymentMethod,
		&tx.Status,
		&tx.Amount,
		&tx.PaidAt,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan transaction", err)
	}

	return &tx, nil
}
