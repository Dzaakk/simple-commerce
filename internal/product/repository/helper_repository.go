package repository

import (
	"Dzaakk/simple-commerce/internal/product/model"
	"Dzaakk/simple-commerce/internal/shopping_cart/models"
	"Dzaakk/simple-commerce/package/template"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func verifyStockAvailability(tx *sql.Tx, listItem []*models.TCartItemDetail) ([]*int, error) {
	var query strings.Builder
	var args []interface{}
	listEmptyProductId := []*int{}

	query.WriteString(`SELECT id, stock FROM public.product WHERE id IN (`)

	for i, item := range listItem {
		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString(fmt.Sprintf("$%d", i+1))
		args = append(args, item.ProductID)
	}

	query.WriteString(") FOR UPDATE")

	rows, err := tx.Query(query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("error locking products: %w", err)
	}
	defer rows.Close()

	stockMap := make(map[int]int)
	for rows.Next() {
		var productID, stock int
		if err := rows.Scan(&productID, &stock); err != nil {
			return nil, fmt.Errorf("error scanning product stock: %w", err)
		}
		stockMap[productID] = stock
	}

	for _, item := range listItem {
		currentStock, exists := stockMap[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("product %d not found", item.ProductID)
		}
		if currentStock == 1 {
			listEmptyProductId = append(listEmptyProductId, &item.ProductID)
		}
		if currentStock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %d: requested %d, available %d", item.ProductID, item.Quantity, currentStock)
		}
	}
	return listEmptyProductId, nil
}

func generateMultipleStockUpdateQuery(listData []*models.TCartItemDetail) (string, []interface{}) {
	var query strings.Builder
	var args []interface{}
	query.WriteString("UPDATE public.product SET stock = CASE id ")

	for _, item := range listData {
		query.WriteString(fmt.Sprintf("WHEN $%d THEN stock - $%d ", len(args)+1, len(args)+2))
		args = append(args, item.ProductID, item.Quantity)
	}

	query.WriteString(" END, updated_by = 'SYSTEM', updated = NOW() WHERE id IN (")

	for i, item := range listData {
		if i > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("$%d", len(args)+1))
		args = append(args, item.ProductID)
	}

	query.WriteString(")")
	return query.String(), args
}

func rowsToProduct(rows *sql.Rows) (*model.TProduct, error) {
	base := template.Base{}
	product := model.TProduct{}

	err := rows.Scan(&product.ID, &product.ProductName, &product.Price, &product.Stock, &product.CategoryID, &base.Created, &base.CreatedBy, &base.Updated, &base.UpdatedBy)
	if err != nil {
		return nil, err
	}
	if !base.UpdatedBy.Valid {
		base.UpdatedBy.String = ""
	}
	product.Base = base

	return &product, nil
}

func retrieveProduct(rows *sql.Rows) (*model.TProduct, error) {
	if rows.Next() {
		return rowsToProduct(rows)
	}
	return nil, errors.New("product not found")
}

func scanProducts(rows *sql.Rows) ([]*model.TProduct, error) {
	var products []*model.TProduct

	for rows.Next() {
		product := &model.TProduct{}
		base := template.Base{}
		var updated sql.NullTime

		err := rows.Scan(
			&product.ID, &product.ProductName, &product.Price, &product.Stock, &product.CategoryID,
			&base.Created, &base.CreatedBy, &updated, &base.UpdatedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		if updated.Valid {
			base.Updated.Time = updated.Time
		}
		if !base.UpdatedBy.Valid {
			base.UpdatedBy.String = ""
		}

		product.Base = base
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return products, nil
}
