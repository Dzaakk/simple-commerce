//go:build wireinject
// +build wireinject

package customer

import (
	handler "Dzaakk/simple-commerce/internal/customer/handler"
	repository "Dzaakk/simple-commerce/internal/customer/repository"
	usecase "Dzaakk/simple-commerce/internal/customer/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *handler.CustomerHandler {
	wire.Build(
		repository.NewCustomerRepository,
		usecase.NewCustomerUseCase,
		handler.NewCustomerHandler,
	)

	return &handler.CustomerHandler{}
}
