//go:build wireinject
// +build wireinject

package injector

import (
	handler "Dzaakk/simple-commerce/internal/product/handlers"
	repository "Dzaakk/simple-commerce/internal/product/repositories"
	routes "Dzaakk/simple-commerce/internal/product/routes"
	usecase "Dzaakk/simple-commerce/internal/product/usecases"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *routes.ProductRoutes {
	wire.Build(
		repository.NewProductRepository,
		usecase.NewProductUseCase,
		handler.NewProductHandler,
		routes.NewProductRoutes,
	)

	return &routes.ProductRoutes{}
}
