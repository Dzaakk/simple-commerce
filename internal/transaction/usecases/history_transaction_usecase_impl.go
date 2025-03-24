package usecases

import (
	"Dzaakk/simple-commerce/internal/transaction/models"
	model "Dzaakk/simple-commerce/internal/transaction/models"
	repo "Dzaakk/simple-commerce/internal/transaction/repositories"
	"context"
	"fmt"
)

type HistoryTransactionUseCaseImpl struct {
	repo repo.HistoryTransactionRepository
}

func NewHistoryTransactionUseCase(repo repo.HistoryTransactionRepository) HistoryTransactionUseCase {
	return &HistoryTransactionUseCaseImpl{repo}
}

func (t *HistoryTransactionUseCaseImpl) CreateHistoryTransaction(ctx context.Context, transactionItem models.THistoryTransaction) {
	panic("unimplemented")
}

func (t *HistoryTransactionUseCaseImpl) GetHistoryTransactionDetail(ctx context.Context, transactionId int64) {
	panic("unimplemented")
}

func (t *HistoryTransactionUseCaseImpl) GetListHistoryTransaction(ctx context.Context, customerId int64) ([]*model.HistoryTransaction, error) {
	listData, err := t.repo.FindByCustomerID(ctx, customerId)
	if err != nil {
		return nil, err
	}

	var result []*model.HistoryTransaction

	for _, d := range listData {
		data := &model.HistoryTransaction{
			TransactionId: fmt.Sprintf("%d", d.TransactionId),
			CustomerId:    fmt.Sprintf("%d", d.CustomerId),
			ProductName:   d.ProductName,
			Price:         fmt.Sprintf("%.2f", d.Price),
			Quantity:      fmt.Sprintf("%d", d.Quantity),
			Status:        d.Status,
		}
		result = append(result, data)
	}

	return result, nil
}
