package service

import (
	catalogDto "Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/internal/order/dto"
	"Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/package/constant"
	"context"
	"database/sql"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderReq) (*dto.OrderRes, error)
	GetOrderByID(ctx context.Context, customerID string, orderID string) (*dto.OrderDetailRes, error)
	GetOrdersByCustomer(ctx context.Context, customerID string, filter dto.OrderFilter) ([]*dto.OrderRes, error)
	CancelOrder(ctx context.Context, customerID string, orderID string) error
	UpdateOrderStatus(ctx context.Context, tx *sql.Tx, orderID string, status constant.OrderStatus) error
}

type OrderRepository interface {
	Create(ctx context.Context, tx *sql.Tx, data *model.Order) (string, error)
	FindByID(ctx context.Context, orderID string) (*model.Order, error)
	FindByCustomerID(ctx context.Context, customerID string, filter dto.OrderFilter) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, tx *sql.Tx, orderID string, status constant.OrderStatus) error
	GenerateOrderNumber(ctx context.Context) (string, error)
}

type OrderItemRepository interface {
	CreateBatch(ctx context.Context, tx *sql.Tx, items []*model.OrderItem) error
	FindByOrderID(ctx context.Context, orderID string) ([]*model.OrderItem, error)
}

type ProductService interface {
	FindByID(ctx context.Context, productID string) (*catalogDto.ProductRes, error)
}

type InventoryService interface {
	ReserveStock(ctx context.Context, tx *sql.Tx, productID string, qty int) error
	ReleaseStock(ctx context.Context, tx *sql.Tx, productID string, qty int) error
}
