//go:build wireinject
// +build wireinject

package product

import (
	handler "Dzaakk/simple-commerce/internal/product/handler"
	repository "Dzaakk/simple-commerce/internal/product/repository"
	usecase "Dzaakk/simple-commerce/internal/product/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *handler.ProductHandler {
	wire.Build(
		repository.NewProductRepository,
		usecase.NewProductUseCase,
		handler.NewProductHandler,
	)

	return &handler.ProductHandler{}
}
