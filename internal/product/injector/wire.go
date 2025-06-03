//go:build wireinject
// +build wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/product/handler"
	"Dzaakk/simple-commerce/internal/product/repository"
	"Dzaakk/simple-commerce/internal/product/route"
	"Dzaakk/simple-commerce/internal/product/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *route.ProductRoutes {
	wire.Build(
		repository.NewProductRepository,
		usecase.NewProductUseCase,
		handler.NewProductHandler,
		route.NewProductRoutes,
	)

	return &route.ProductRoutes{}
}
