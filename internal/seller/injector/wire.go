//go:build wireinject
// +build wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/seller/handler"
	"Dzaakk/simple-commerce/internal/seller/repository"
	"Dzaakk/simple-commerce/internal/seller/route"
	"Dzaakk/simple-commerce/internal/seller/usecase"

	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *route.SellerRoutes {
	wire.Build(
		repository.NewSellerRepository,
		usecase.NewSellerUseCase,
		handler.NewSellerHandler,
		route.NewSellerRoutes,
	)

	return &route.SellerRoutes{}
}
