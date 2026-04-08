package route

import (
	"Dzaakk/simple-commerce/internal/order/handler"
	"Dzaakk/simple-commerce/internal/order/repository"
	"Dzaakk/simple-commerce/internal/order/service"
	catalogRepo "Dzaakk/simple-commerce/internal/catalog/repository"
	catalogService "Dzaakk/simple-commerce/internal/catalog/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type OrderRoutes struct {
	Handler *handler.OrderHandler
}

func NewOrderRoutes(handler *handler.OrderHandler) *OrderRoutes {
	return &OrderRoutes{Handler: handler}
}

func (or *OrderRoutes) Route(r *gin.RouterGroup) {
	orders := r.Group("/api/v1/orders")
	{
		orders.POST("", or.Handler.CreateOrder)
		orders.GET("", or.Handler.GetOrders)
		orders.GET("/:id", or.Handler.GetOrderDetail)
		orders.PATCH("/:id/cancel", or.Handler.CancelOrder)
	}
}

func InitializedService(db *sql.DB) *OrderRoutes {
	orderRepo := repository.NewOrderRepository(db)
	orderItemRepo := repository.NewOrderItemRepository(db)

	productRepo := catalogRepo.NewProductRepository(db)
	inventoryRepo := catalogRepo.NewInventoryRepository(db)

	productService := catalogService.NewProductService(productRepo)
	inventoryService := catalogService.NewInventoryService(inventoryRepo)

	orderService := service.NewOrderService(db, orderRepo, orderItemRepo, productService, inventoryService)
	orderHandler := handler.NewOrderHandler(orderService)

	return NewOrderRoutes(orderHandler)
}
