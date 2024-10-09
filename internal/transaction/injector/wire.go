//go:build wireinject
// +build wireinject

package transaction

import (
	customer "Dzaakk/simple-commerce/internal/customer/repository"
	cart "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	handler "Dzaakk/simple-commerce/internal/transaction/handler"
	repository "Dzaakk/simple-commerce/internal/transaction/repository"
	routes "Dzaakk/simple-commerce/internal/transaction/routes"
	usecase "Dzaakk/simple-commerce/internal/transaction/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *handler.TransactionHandler {
	wire.Build(
		repository.NewTransactionRepository,
		usecase.NewTransactionUseCase,
		handler.NewTransactionHandler,
		cart.NewShoppingCartRepository,
		cart.NewShoppingCartItemRepository,
		customer.NewCustomerRepository,
		routes.NewTransactionRoutes,
	)

	return &handler.TransactionHandler{}
}
