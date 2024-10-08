// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	repository3 "Dzaakk/simple-commerce/internal/customer/repository"
	repository2 "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/internal/transaction/handler"
	"Dzaakk/simple-commerce/internal/transaction/repository"
	"Dzaakk/simple-commerce/internal/transaction/routes"
	"Dzaakk/simple-commerce/internal/transaction/usecase"
	"database/sql"
)

// Injectors from wire.go:

func InitializedService(db *sql.DB) *routes.TransactionRoutes {
	transactionRepository := repository.NewTransactionRepository(db)
	shoppingCartRepository := repository2.NewShoppingCartRepository(db)
	shoppingCartItemRepository := repository2.NewShoppingCartItemRepository(db)
	customerRepository := repository3.NewCustomerRepository(db)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepository, shoppingCartRepository, shoppingCartItemRepository, customerRepository)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	transactionRoutes := routes.NewTransactionRoutes(transactionHandler)
	return transactionRoutes
}
