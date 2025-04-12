package usecases

import (
	modelItem "Dzaakk/simple-commerce/internal/shopping_cart/models"
	model "Dzaakk/simple-commerce/internal/transaction/models"
	template "Dzaakk/simple-commerce/package/templates"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func insertToTableTransactionWithTx(ctx context.Context, tx *sql.Tx, t *TransactionUseCaseImpl, customerID, cartID, totalTransaction int) (*string, error) {
	newTransaction := model.TTransaction{
		CustomerID:      customerID,
		CartID:          cartID,
		TotalAmount:     float32(totalTransaction),
		TransactionDate: time.Now(),
		Status:          "Success",
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: fmt.Sprintf("%d", customerID),
		},
	}

	data, err := t.repo.CreateWithTx(ctx, tx, newTransaction)
	if err != nil {
		return nil, err
	}
	transactionDate := data.TransactionDate.Format("06-01-02 15:04:05")

	return &transactionDate, nil
}

func generateReceipt(listItem []*modelItem.TCartItemDetail) (*model.TransactionRes, error) {
	var res model.TransactionRes
	var listProduct []model.ListProduct
	total := 0
	for _, item := range listItem {
		product := model.ListProduct{
			ProductName: item.ProductName,
			Price:       fmt.Sprintf("%.0f", item.Price),
			Quantity:    fmt.Sprintf("%d", item.Quantity),
		}
		listProduct = append(listProduct, product)
		total = total + (int(item.Price) * item.Quantity)
	}
	res.ListProduct = listProduct
	res.TotalTransaction = fmt.Sprintf("%d", total)

	return &res, nil
}
