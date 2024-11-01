package repository

import (
	"Dzaakk/simple-commerce/internal/shopping_cart/models"
	"database/sql"
	"fmt"
	"strings"
)

type TransactionItemRepositoryImpl struct {
	DB *sql.DB
}

func NewTransactionItemRepository(db *sql.DB) TransactionItemRepository {
	return &TransactionItemRepositoryImpl{
		DB: db,
	}
}

func (t *TransactionItemRepositoryImpl) Create(data []*models.TCartItemDetail, customerId int64) error {
	if len(data) == 0 {
		return nil
	}

	listQuery := generateInsertStatements(data, customerId)

	tx, err := t.DB.Begin()
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
		_, err := tx.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute insert: %v, error: %w", query, err)
		}
	}

	return nil
}

func generateInsertStatements(listData []*models.TCartItemDetail, customerId int64) []string {
	var sqlInserts []string
	columns := "customer_id, productName, price, quantity, status"
	for _, data := range listData {
		values := []interface{}{
			customerId, data.ProductName, data.Price, data.Quantity, "PAID",
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
