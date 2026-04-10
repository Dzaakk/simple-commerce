package service

import (
	orderModel "Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/internal/transaction/dto"
	txModel "Dzaakk/simple-commerce/internal/transaction/model"
	"Dzaakk/simple-commerce/package/constant"
	"context"
	"database/sql"
	"time"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionReq) (*dto.TransactionRes, error)
	GetTransactionByID(ctx context.Context, customerID, transactionID string) (*dto.TransactionRes, error)
	GetTransactionByOrderID(ctx context.Context, customerID, orderID string) (*dto.TransactionRes, error)
	HandlePaymentCallback(ctx context.Context, req *dto.PaymentCallbackReq) error
	ExpireTransaction(ctx context.Context, transactionID string) error
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, data *txModel.Transaction) (string, error)
	FindByID(ctx context.Context, transactionID string) (*txModel.Transaction, error)
	FindByOrderID(ctx context.Context, orderID string) (*txModel.Transaction, error)
	FindByTransactionNumber(ctx context.Context, txNumber string) (*txModel.Transaction, error)
	UpdateStatus(ctx context.Context, tx *sql.Tx, transactionID string, status constant.TransactionStatus, paidAt *time.Time) error
	GenerateTransactionNumber(ctx context.Context) (string, error)
}

type OrderRepository interface {
	FindByID(ctx context.Context, orderID string) (*orderModel.Order, error)
}

type OrderItemRepository interface {
	FindByOrderID(ctx context.Context, orderID string) ([]*orderModel.OrderItem, error)
}

type OrderService interface {
	UpdateOrderStatus(ctx context.Context, tx *sql.Tx, orderID string, status constant.OrderStatus) error
}

type InventoryService interface {
	ReleaseStock(ctx context.Context, tx *sql.Tx, productID string, qty int) error
}
