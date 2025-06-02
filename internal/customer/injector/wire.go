//go:build wireinject
// +build wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/customer/handler"
	"Dzaakk/simple-commerce/internal/customer/repository"
	"Dzaakk/simple-commerce/internal/customer/route"
	"Dzaakk/simple-commerce/internal/customer/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *route.CustomerRoutes {
	wire.Build(
		repository.NewCustomerRepository,
		usecase.NewCustomerUseCase,
		handler.NewCustomerHandler,
		route.NewCustomerRoutes,
	)

	return &route.CustomerRoutes{}
}
