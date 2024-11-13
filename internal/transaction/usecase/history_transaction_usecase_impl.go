package usecase

import (
	"Dzaakk/simple-commerce/internal/transaction/models"
	model "Dzaakk/simple-commerce/internal/transaction/models"
	repo "Dzaakk/simple-commerce/internal/transaction/repository"
	"fmt"
)

type HistoryTransactionUseCaseImpl struct {
	repo repo.HistoryTransactionRepository
}

func NewHistoryTransactionUseCase(repo repo.HistoryTransactionRepository) HistoryTransactionUseCase {
	return &HistoryTransactionUseCaseImpl{repo}
}

func (t *HistoryTransactionUseCaseImpl) CreateHistoryTransaction(transactionItem models.THistoryTransaction) {
	panic("unimplemented")
}

func (t *HistoryTransactionUseCaseImpl) GetHistoryTransactionDetail(transactionId int64) {
	panic("unimplemented")
}

func (t *HistoryTransactionUseCaseImpl) GetListHistoryTransaction(customerId int64) ([]*model.HistoryTransaction, error) {
	listData, err := t.repo.FindByCustomerId(customerId)
	if err != nil {
		return nil, err
	}

	var result []*model.HistoryTransaction

	for _, d := range listData {
		data := &model.HistoryTransaction{
			Id:          fmt.Sprintf("%d", d.Id),
			CustomerId:  fmt.Sprintf("%d", d.CustomerId),
			ProductName: d.ProductName,
			Price:       fmt.Sprintf("%.2f", d.Price),
			Quantity:    fmt.Sprintf("%d", d.Quantity),
			Status:      d.Status,
		}
		result = append(result, data)
	}

	return result, nil
}
