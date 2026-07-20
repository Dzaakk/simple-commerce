package service

import (
	"context"
	"errors"
	"testing"
	"time"

	orderModel "Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/internal/transaction/dto"
	txModel "Dzaakk/simple-commerce/internal/transaction/model"
	"Dzaakk/simple-commerce/package/constant"
)

type mockTransactionTransactor struct {
	called bool
}

func (m *mockTransactionTransactor) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	m.called = true
	return fn(ctx)
}

type mockTransactionRepository struct {
	createFn                  func(context.Context, *txModel.Transaction) (string, error)
	findByIDFn                func(context.Context, string) (*txModel.Transaction, error)
	findByOrderIDFn           func(context.Context, string) (*txModel.Transaction, error)
	findByTransactionNumberFn func(context.Context, string) (*txModel.Transaction, error)
	updateStatusFn            func(context.Context, string, constant.TransactionStatus, *time.Time) error
	generateTransactionNumber func(context.Context) (string, error)
}

func (m *mockTransactionRepository) Create(ctx context.Context, data *txModel.Transaction) (string, error) {
	if m.createFn == nil {
		return "", errors.New("unexpected Create call")
	}
	return m.createFn(ctx, data)
}

func (m *mockTransactionRepository) FindByID(ctx context.Context, transactionID string) (*txModel.Transaction, error) {
	if m.findByIDFn == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return m.findByIDFn(ctx, transactionID)
}

func (m *mockTransactionRepository) FindByOrderID(ctx context.Context, orderID string) (*txModel.Transaction, error) {
	if m.findByOrderIDFn == nil {
		return nil, errors.New("unexpected FindByOrderID call")
	}
	return m.findByOrderIDFn(ctx, orderID)
}

func (m *mockTransactionRepository) FindByTransactionNumber(ctx context.Context, txNumber string) (*txModel.Transaction, error) {
	if m.findByTransactionNumberFn == nil {
		return nil, errors.New("unexpected FindByTransactionNumber call")
	}
	return m.findByTransactionNumberFn(ctx, txNumber)
}

func (m *mockTransactionRepository) UpdateStatus(ctx context.Context, transactionID string, status constant.TransactionStatus, paidAt *time.Time) error {
	if m.updateStatusFn == nil {
		return errors.New("unexpected UpdateStatus call")
	}
	return m.updateStatusFn(ctx, transactionID, status, paidAt)
}

func (m *mockTransactionRepository) GenerateTransactionNumber(ctx context.Context) (string, error) {
	if m.generateTransactionNumber == nil {
		return "", errors.New("unexpected GenerateTransactionNumber call")
	}
	return m.generateTransactionNumber(ctx)
}

type mockTransactionOrderRepository struct {
	findByIDFn func(context.Context, string) (*orderModel.Order, error)
}

func (m *mockTransactionOrderRepository) FindByID(ctx context.Context, orderID string) (*orderModel.Order, error) {
	if m.findByIDFn == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return m.findByIDFn(ctx, orderID)
}

type mockTransactionOrderItemRepository struct {
	findByOrderIDFn func(context.Context, string) ([]*orderModel.OrderItem, error)
}

func (m *mockTransactionOrderItemRepository) FindByOrderID(ctx context.Context, orderID string) ([]*orderModel.OrderItem, error) {
	if m.findByOrderIDFn == nil {
		return nil, errors.New("unexpected FindByOrderID call")
	}
	return m.findByOrderIDFn(ctx, orderID)
}

type mockTransactionOrderService struct {
	updateOrderStatusFn func(context.Context, string, constant.OrderStatus) error
}

func (m *mockTransactionOrderService) UpdateOrderStatus(ctx context.Context, orderID string, status constant.OrderStatus) error {
	if m.updateOrderStatusFn == nil {
		return errors.New("unexpected UpdateOrderStatus call")
	}
	return m.updateOrderStatusFn(ctx, orderID, status)
}

type mockTransactionInventoryService struct {
	releaseStockFn func(context.Context, string, int) error
}

func (m *mockTransactionInventoryService) ReleaseStock(ctx context.Context, productID string, qty int) error {
	if m.releaseStockFn == nil {
		return errors.New("unexpected ReleaseStock call")
	}
	return m.releaseStockFn(ctx, productID, qty)
}

func TestTransactionServiceHandlePaymentCallbackRunsUpdatesWithinTransaction(t *testing.T) {
	txManager := &mockTransactionTransactor{}
	txRepo := &mockTransactionRepository{
		findByTransactionNumberFn: func(_ context.Context, txNumber string) (*txModel.Transaction, error) {
			if txNumber != "TRX-1" {
				t.Fatalf("transaction number = %q, want TRX-1", txNumber)
			}
			return &txModel.Transaction{ID: "transaction-1", OrderID: "order-1", Status: string(constant.TransactionPending)}, nil
		},
		updateStatusFn: func(_ context.Context, transactionID string, status constant.TransactionStatus, paidAt *time.Time) error {
			if transactionID != "transaction-1" {
				t.Fatalf("transaction id = %q, want transaction-1", transactionID)
			}
			if status != constant.TransactionSuccess {
				t.Fatalf("status = %q, want %q", status, constant.TransactionSuccess)
			}
			if paidAt == nil {
				t.Fatal("paidAt must be set for successful payment callback")
			}
			return nil
		},
	}
	orderRepo := &mockTransactionOrderRepository{
		findByIDFn: func(_ context.Context, orderID string) (*orderModel.Order, error) {
			if orderID != "order-1" {
				t.Fatalf("order id = %q, want order-1", orderID)
			}
			return &orderModel.Order{ID: orderID, CustomerID: "customer-1", Status: string(constant.OrderPending)}, nil
		},
	}
	orderItemRepo := &mockTransactionOrderItemRepository{
		findByOrderIDFn: func(context.Context, string) ([]*orderModel.OrderItem, error) {
			return nil, nil
		},
	}
	orderSvc := &mockTransactionOrderService{
		updateOrderStatusFn: func(_ context.Context, orderID string, status constant.OrderStatus) error {
			if orderID != "order-1" {
				t.Fatalf("order id = %q, want order-1", orderID)
			}
			if status != constant.OrderConfirmed {
				t.Fatalf("order status = %q, want %q", status, constant.OrderConfirmed)
			}
			return nil
		},
	}

	err := NewTransactionService(txManager, txRepo, orderRepo, orderItemRepo, orderSvc, &mockTransactionInventoryService{}).
		HandlePaymentCallback(context.Background(), &dto.PaymentCallbackReq{
			TransactionNumber: "TRX-1",
			Status:            constant.TransactionSuccess,
			Signature:         "valid-signature",
		})
	if err != nil {
		t.Fatalf("HandlePaymentCallback returned error: %v", err)
	}
	if !txManager.called {
		t.Fatal("HandlePaymentCallback must run updates inside a transaction")
	}
}
