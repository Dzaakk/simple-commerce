//go:build wireinject
// +build wireinject

package customer

import (
	handler "Dzaakk/synapsis/internal/customer/handler"
	repository "Dzaakk/synapsis/internal/customer/repository"
	usecase "Dzaakk/synapsis/internal/customer/usecase"
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
