//go:build wireinject
// +build wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/seller/handlers"
	"Dzaakk/simple-commerce/internal/seller/repositories"
	"Dzaakk/simple-commerce/internal/seller/routes"
	"Dzaakk/simple-commerce/internal/seller/usecases"

	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *routes.SellerRoutes {
	wire.Build(
		repositories.NewSellerRepository,
		usecases.NewSellerUseCase,
		handlers.NewSellerHandler,
		routes.NewSellerRoutes,
	)

	return &routes.SellerRoutes{}
}
