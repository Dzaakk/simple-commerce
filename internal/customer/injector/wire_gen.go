// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/customer/handlers"
	"Dzaakk/simple-commerce/internal/customer/repositories"
	"Dzaakk/simple-commerce/internal/customer/routes"
	"Dzaakk/simple-commerce/internal/customer/usecases"
	"database/sql"
)

// Injectors from wire.go:

func InitializedService(db *sql.DB) *routes.CustomerRoutes {
	customerRepository := repository.NewCustomerRepository(db)
	customerUseCase := usecase.NewCustomerUseCase(customerRepository)
	customerHandler := handler.NewCustomerHandler(customerUseCase)
	customerRoutes := routes.NewCustomerRoutes(customerHandler)
	return customerRoutes
}
