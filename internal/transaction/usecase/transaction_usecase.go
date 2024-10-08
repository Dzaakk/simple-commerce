package transaction

import (
	customer "Dzaakk/simple-commerce/internal/customer/repository"
	modelItem "Dzaakk/simple-commerce/internal/shopping_cart/models"
	shoppingCart "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	model "Dzaakk/simple-commerce/internal/transaction/models"
	repo "Dzaakk/simple-commerce/internal/transaction/repository"
	"Dzaakk/simple-commerce/package/template"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type TransactionUseCase interface {
	CreateTransaction(data model.TransactionReq) (*model.TransactionRes, error)
}

type TransactionUseCaseImpl struct {
	repo         repo.TransactionRepository
	repoCart     shoppingCart.ShoppingCartRepository
	repoCartItem shoppingCart.ShoppingCartItemRepository
	repoCustomer customer.CustomerRepository
}

func NewTransactionUseCase(repo repo.TransactionRepository, repoCart shoppingCart.ShoppingCartRepository, repoCartItem shoppingCart.ShoppingCartItemRepository, repoCustomer customer.CustomerRepository) TransactionUseCase {
	return &TransactionUseCaseImpl{repo, repoCart, repoCartItem, repoCustomer}
}

func (t *TransactionUseCaseImpl) CreateTransaction(data model.TransactionReq) (*model.TransactionRes, error) {
	cartId, _ := strconv.Atoi(data.CartId)
	customerId, _ := strconv.Atoi(data.CustomerId)

	listItem, err := t.repoCartItem.RetrieveCartItemsByCartId(cartId) // get all items on cart
	if err != nil {
		return nil, err
	}
	res, err := generateReceipt(listItem) // generate receipt and calculate total transaction
	if err != nil {
		return nil, err
	}

	customer, err := t.repoCustomer.GetBalance(customerId) // check customer current balance
	if err != nil {
		return nil, err
	}
	totalTransaction, _ := strconv.Atoi(res.TotalTransaction)
	if totalTransaction > int(customer.Balance) {
		return nil, errors.New("insufficient balance")
	}

	err = t.repoCartItem.DeleteAll(cartId) // delete all item on cart
	if err != nil {
		return nil, err
	}

	_, err = t.repoCart.UpdateStatusById(cartId, "Paid", data.CustomerId) // update cart status to 'Paid'
	if err != nil {
		return nil, err
	}

	transactionDate, err := insertToTableTransaction(t, customerId, cartId, totalTransaction)
	if err != nil {
		return nil, err
	}

	newBalance := customer.Balance - float32(totalTransaction)
	_, err = t.repoCustomer.UpdateBalance(customerId, newBalance) // update balance customer
	if err != nil {
		return nil, err
	}

	res.CustomerId = data.CustomerId
	res.TransactionDate = *transactionDate
	return res, nil
}

func insertToTableTransaction(t *TransactionUseCaseImpl, customerId, cartId, totalTransaction int) (*string, error) {
	newTransaction := model.TTransaction{
		CustomerId:      customerId,
		CartId:          cartId,
		TotalAmount:     float32(totalTransaction),
		TransactionDate: time.Now(),
		Status:          "Success",
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: fmt.Sprintf("%d", customerId),
		},
	}

	data, err := t.repo.Create(newTransaction)
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
