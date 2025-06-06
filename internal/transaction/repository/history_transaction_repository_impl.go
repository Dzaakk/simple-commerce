package repository

import (
	cartModel "Dzaakk/simple-commerce/internal/shopping_cart/model"
	"Dzaakk/simple-commerce/internal/transaction/model"
	"Dzaakk/simple-commerce/package/template"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type HistoryTransactionRepositoryImpl struct {
	DB *sql.DB
}

func NewHistoryTransactionRepository(db *sql.DB) HistoryTransactionRepository {
	return &HistoryTransactionRepositoryImpl{
		DB: db,
	}
}

func (repo *HistoryTransactionRepositoryImpl) Create(ctx context.Context, data []*cartModel.TCartItemDetail, customerID int64) error {
	if len(data) == 0 {
		return nil
	}

	listQuery := generateInsertStatements(data, customerID)

	tx, err := repo.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, query := range listQuery {
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute insert: %v, error: %w", query, err)
		}
	}

	return nil
}

const queryFindByCustomerID = "SELECT * FROM public.history_transaction WHERE customer_id=$1"

func (repo *HistoryTransactionRepositoryImpl) FindByCustomerID(ctx context.Context, customerID int64) ([]*model.THistoryTransaction, error) {
	rows, err := repo.DB.QueryContext(ctx, queryFindByCustomerID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listHistoryTransaction []*model.THistoryTransaction
	for rows.Next() {
		historyTransaction, err := retrieveHistoryTransaction(rows)
		if err != nil {
			return nil, err
		}
		listHistoryTransaction = append(listHistoryTransaction, historyTransaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return listHistoryTransaction, nil
}

func generateInsertStatements(listData []*cartModel.TCartItemDetail, customerID int64) []string {
	var sqlInserts []string
	columns := "customer_id, productName, price, quantity, status"
	for _, data := range listData {
		values := []interface{}{
			customerID, data.ProductName, data.Price, data.Quantity, "PAID",
		}
		sqlInsert := fmt.Sprintf("INSERT INTO transaction_items (%s) VALUES (%s);",
			columns, formatValues(values))
		sqlInserts = append(sqlInserts, sqlInsert)
	}
	return sqlInserts
}

func formatValues(values []interface{}) string {
	var formattedValues []string
	for _, v := range values {
		switch v := v.(type) {
		case string:
			formattedValues = append(formattedValues, fmt.Sprintf("'%s'", v))
		case float64:
			formattedValues = append(formattedValues, fmt.Sprintf("'%.2f'", v))
		default:
			formattedValues = append(formattedValues, fmt.Sprintf("'%v'", v))
		}
	}
	return strings.Join(formattedValues, ", ")
}

func retrieveHistoryTransaction(rows *sql.Rows) (*model.THistoryTransaction, error) {
	if rows.Next() {
		return rowsToProduct(rows)
	}
	return nil, errors.New("record not found")
}

func rowsToProduct(rows *sql.Rows) (*model.THistoryTransaction, error) {
	base := template.Base{}
	ht := model.THistoryTransaction{}

	err := rows.Scan(ht.TransactionID, ht.CustomerID, ht.Price, ht.ProductName, ht.Status, ht.Quantity, base.Created, base.CreatedBy, base.Updated, base.UpdatedBy)
	if err != nil {
		return nil, err
	}
	if !base.UpdatedBy.Valid {
		base.UpdatedBy.String = ""
	}
	ht.Base = base

	return &ht, nil
}
