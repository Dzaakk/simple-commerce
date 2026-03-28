package route

import (
	"Dzaakk/simple-commerce/internal/user/handler"
	"Dzaakk/simple-commerce/internal/user/repository"
	"Dzaakk/simple-commerce/internal/user/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	Handler *handler.UserHandler
}

func NewuserRoutes(handler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		Handler: handler,
	}
}

func (ur *UserRoutes) Route(r *gin.RouterGroup) {
	userHandler := r.Group("/api/v1/customer")
	userHandler.Use()
	{
		userHandler.GET("/by-email", ur.Handler.FindCustomerByEmail)
		userHandler.GET("/by-id", ur.Handler.FindCustomerByID)
	}

	sellerHandler := r.Group("/api/v1/seller")
	{
		sellerHandler.GET("/by-name", ur.Handler.FindSellerByName)
	}
}

func InitializedService(db *sql.DB) *UserRoutes {
	customerRepo := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepo)
	sellerRepo := repository.NewSellerRepository(db)
	sellerService := service.NewSellerService(sellerRepo)
	handler := handler.NewUserHandler(customerService, sellerService)

	return NewuserRoutes(handler)
}
