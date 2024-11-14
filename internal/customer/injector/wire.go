//go:build wireinject
// +build wireinject

package injector

import (
	handlers "Dzaakk/simple-commerce/internal/customer/handlers"
	repositories "Dzaakk/simple-commerce/internal/customer/repositories"
	routes "Dzaakk/simple-commerce/internal/customer/routes"
	usecases "Dzaakk/simple-commerce/internal/customer/usecases"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *routes.CustomerRoutes {
	wire.Build(
		repositories.NewCustomerRepository,
		usecases.NewCustomerUseCase,
		handlers.NewCustomerHandler,
		routes.NewCustomerRoutes,
	)

	return &routes.CustomerRoutes{}
}
