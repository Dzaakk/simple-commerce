//go:build wireinject
// +build wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/product/handlers"
	"Dzaakk/simple-commerce/internal/product/repositories"
	"Dzaakk/simple-commerce/internal/product/routes"
	"Dzaakk/simple-commerce/internal/product/usecases"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *routes.ProductRoutes {
	wire.Build(
		repositories.NewProductRepository,
		usecases.NewProductUseCase,
		handlers.NewProductHandler,
		routes.NewProductRoutes,
	)

	return &routes.ProductRoutes{}
}
