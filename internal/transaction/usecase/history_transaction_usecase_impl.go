package usecase

import (
	"Dzaakk/simple-commerce/internal/transaction/model"
	repo "Dzaakk/simple-commerce/internal/transaction/repository"
	"context"
	"fmt"
)

type HistoryTransactionUseCaseImpl struct {
	repo repo.HistoryTransactionRepository
}

func NewHistoryTransactionUseCase(repo repo.HistoryTransactionRepository) HistoryTransactionUseCase {
	return &HistoryTransactionUseCaseImpl{repo}
}

func (t *HistoryTransactionUseCaseImpl) CreateHistoryTransaction(ctx context.Context, transactionItem model.THistoryTransaction) {
	panic("unimplemented")
}

func (t *HistoryTransactionUseCaseImpl) GetHistoryTransactionDetail(ctx context.Context, transactionID int64) {
	panic("unimplemented")
}

func (t *HistoryTransactionUseCaseImpl) GetListHistoryTransaction(ctx context.Context, customerID int64) ([]*model.HistoryTransaction, error) {
	listData, err := t.repo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	var result []*model.HistoryTransaction

	for _, d := range listData {
		data := &model.HistoryTransaction{
			TransactionID: fmt.Sprintf("%d", d.TransactionID),
			CustomerID:    fmt.Sprintf("%d", d.CustomerID),
			ProductName:   d.ProductName,
			Price:         fmt.Sprintf("%.2f", d.Price),
			Quantity:      fmt.Sprintf("%d", d.Quantity),
			Status:        d.Status,
		}
		result = append(result, data)
	}

	return result, nil
}
