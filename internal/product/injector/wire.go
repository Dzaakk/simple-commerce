//go:build wireinject
// +build wireinject

package product

import (
	handler "Dzaakk/synapsis/internal/product/handler"
	repository "Dzaakk/synapsis/internal/product/repository"
	usecase "Dzaakk/synapsis/internal/product/usecase"
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
