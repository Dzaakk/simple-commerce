//go:build wireinject
// +build wireinject

package injector

import (
	handler "Dzaakk/simple-commerce/internal/customer/handler"
	repository "Dzaakk/simple-commerce/internal/customer/repository"
	routes "Dzaakk/simple-commerce/internal/customer/routes"
	usecase "Dzaakk/simple-commerce/internal/customer/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *routes.CustomerRoutes {
	wire.Build(
		repository.NewCustomerRepository,
		usecase.NewCustomerUseCase,
		handler.NewCustomerHandler,
		routes.NewCustomerRoutes,
	)

	return &routes.CustomerRoutes{}
}
