// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/product/handler"
	"Dzaakk/simple-commerce/internal/product/repository"
	"Dzaakk/simple-commerce/internal/product/route"
	"Dzaakk/simple-commerce/internal/product/usecase"
	"database/sql"
)

// Injectors from wire.go:

func InitializedService(db *sql.DB) *route.ProductRoutes {
	productRepository := repository.NewProductRepository(db)
	productUseCase := usecase.NewProductUseCase(productRepository)
	productHandler := handler.NewProductHandler(productUseCase)
	productRoutes := route.NewProductRoutes(productHandler)
	return productRoutes
}
