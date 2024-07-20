//go:build wireinject
// +build wireinject

package transaction

import (
	customer "Dzaakk/synapsis/internal/customer/repository"
	cart "Dzaakk/synapsis/internal/shopping_cart/repository"
	handler "Dzaakk/synapsis/internal/transaction/handler"
	repository "Dzaakk/synapsis/internal/transaction/repository"
	usecase "Dzaakk/synapsis/internal/transaction/usecase"
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
	)

	return &handler.TransactionHandler{}
}
