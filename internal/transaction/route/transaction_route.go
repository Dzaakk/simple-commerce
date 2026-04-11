package route

import (
	"Dzaakk/simple-commerce/internal/transaction/handler"
	"Dzaakk/simple-commerce/internal/transaction/repository"
	"Dzaakk/simple-commerce/internal/transaction/service"
	orderRepo "Dzaakk/simple-commerce/internal/order/repository"
	orderService "Dzaakk/simple-commerce/internal/order/service"
	catalogRepo "Dzaakk/simple-commerce/internal/catalog/repository"
	catalogService "Dzaakk/simple-commerce/internal/catalog/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type TransactionRoutes struct {
	Handler *handler.TransactionHandler
}

func NewTransactionRoutes(handler *handler.TransactionHandler) *TransactionRoutes {
	return &TransactionRoutes{Handler: handler}
}

func (tr *TransactionRoutes) Route(r *gin.RouterGroup) {
	api := r.Group("/api/v1/transactions")
	{
		api.POST("", tr.Handler.CreateTransaction)
		api.GET("/:id", tr.Handler.GetTransactionByID)
		api.GET("/by-order/:order_id", tr.Handler.GetTransactionByOrderID)
		api.POST("/callback", tr.Handler.PaymentCallback)
		api.PATCH("/:id/expire", tr.Handler.ExpireTransaction)
	}
}

func InitializedService(db *sql.DB) *TransactionRoutes {
	txRepo := repository.NewTransactionRepository(db)
	orderRepoImpl := orderRepo.NewOrderRepository(db)
	orderItemRepo := orderRepo.NewOrderItemRepository(db)

	productRepo := catalogRepo.NewProductRepository(db)
	inventoryRepo := catalogRepo.NewInventoryRepository(db)

	productSvc := catalogService.NewProductService(productRepo)
	inventorySvc := catalogService.NewInventoryService(inventoryRepo)
	orderSvc := orderService.NewOrderService(db, orderRepoImpl, orderItemRepo, productSvc, inventorySvc)

	svc := service.NewTransactionService(db, txRepo, orderRepoImpl, orderItemRepo, orderSvc, inventorySvc)
	h := handler.NewTransactionHandler(svc)

	return NewTransactionRoutes(h)
}
