package service

import (
	"context"
	"errors"
	"testing"

	catalogDto "Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/internal/order/dto"
	"Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/package/constant"
)

type mockOrderTransactor struct {
	called bool
}

func (m *mockOrderTransactor) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	m.called = true
	return fn(ctx)
}

type mockOrderRepository struct {
	generateOrderNumberFn func(context.Context) (string, error)
	createFn              func(context.Context, *model.Order) (string, error)
	findByIDFn            func(context.Context, string) (*model.Order, error)
	findByCustomerIDFn    func(context.Context, string, dto.OrderFilter) ([]*model.Order, error)
	updateStatusFn        func(context.Context, string, constant.OrderStatus) error
}

func (m *mockOrderRepository) Create(ctx context.Context, data *model.Order) (string, error) {
	if m.createFn == nil {
		return "", errors.New("unexpected Create call")
	}
	return m.createFn(ctx, data)
}

func (m *mockOrderRepository) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	if m.findByIDFn == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return m.findByIDFn(ctx, orderID)
}

func (m *mockOrderRepository) FindByCustomerID(ctx context.Context, customerID string, filter dto.OrderFilter) ([]*model.Order, error) {
	if m.findByCustomerIDFn == nil {
		return nil, errors.New("unexpected FindByCustomerID call")
	}
	return m.findByCustomerIDFn(ctx, customerID, filter)
}

func (m *mockOrderRepository) UpdateStatus(ctx context.Context, orderID string, status constant.OrderStatus) error {
	if m.updateStatusFn == nil {
		return errors.New("unexpected UpdateStatus call")
	}
	return m.updateStatusFn(ctx, orderID, status)
}

func (m *mockOrderRepository) GenerateOrderNumber(ctx context.Context) (string, error) {
	if m.generateOrderNumberFn == nil {
		return "", errors.New("unexpected GenerateOrderNumber call")
	}
	return m.generateOrderNumberFn(ctx)
}

type mockOrderItemRepository struct {
	createBatchFn   func(context.Context, []*model.OrderItem) error
	findByOrderIDFn func(context.Context, string) ([]*model.OrderItem, error)
}

func (m *mockOrderItemRepository) CreateBatch(ctx context.Context, items []*model.OrderItem) error {
	if m.createBatchFn == nil {
		return errors.New("unexpected CreateBatch call")
	}
	return m.createBatchFn(ctx, items)
}

func (m *mockOrderItemRepository) FindByOrderID(ctx context.Context, orderID string) ([]*model.OrderItem, error) {
	if m.findByOrderIDFn == nil {
		return nil, errors.New("unexpected FindByOrderID call")
	}
	return m.findByOrderIDFn(ctx, orderID)
}

type mockOrderProductService struct {
	findByIDFn func(context.Context, string) (*catalogDto.ProductRes, error)
}

func (m *mockOrderProductService) FindByID(ctx context.Context, productID string) (*catalogDto.ProductRes, error) {
	if m.findByIDFn == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return m.findByIDFn(ctx, productID)
}

type mockOrderInventoryService struct {
	reserveStockFn func(context.Context, string, int) error
	releaseStockFn func(context.Context, string, int) error
}

func (m *mockOrderInventoryService) ReserveStock(ctx context.Context, productID string, qty int) error {
	if m.reserveStockFn == nil {
		return errors.New("unexpected ReserveStock call")
	}
	return m.reserveStockFn(ctx, productID, qty)
}

func (m *mockOrderInventoryService) ReleaseStock(ctx context.Context, productID string, qty int) error {
	if m.releaseStockFn == nil {
		return errors.New("unexpected ReleaseStock call")
	}
	return m.releaseStockFn(ctx, productID, qty)
}

func TestOrderServiceCreateOrderRunsWritesWithinTransaction(t *testing.T) {
	txManager := &mockOrderTransactor{}
	orderRepo := &mockOrderRepository{
		generateOrderNumberFn: func(context.Context) (string, error) {
			return "ORD-20260720-0001", nil
		},
		createFn: func(_ context.Context, order *model.Order) (string, error) {
			if order.CustomerID != "customer-1" {
				t.Fatalf("customer id = %q, want customer-1", order.CustomerID)
			}
			if order.TotalAmount != 20 {
				t.Fatalf("total amount = %v, want 20", order.TotalAmount)
			}
			return "order-1", nil
		},
	}
	itemRepo := &mockOrderItemRepository{
		createBatchFn: func(_ context.Context, items []*model.OrderItem) error {
			if len(items) != 1 {
				t.Fatalf("items len = %d, want 1", len(items))
			}
			if items[0].OrderID != "order-1" {
				t.Fatalf("item order id = %q, want order-1", items[0].OrderID)
			}
			return nil
		},
	}
	productSvc := &mockOrderProductService{
		findByIDFn: func(_ context.Context, productID string) (*catalogDto.ProductRes, error) {
			if productID != "product-1" {
				t.Fatalf("product id = %q, want product-1", productID)
			}
			return &catalogDto.ProductRes{ID: productID, SellerID: "seller-1", Price: 10, IsActive: true}, nil
		},
	}
	inventorySvc := &mockOrderInventoryService{
		reserveStockFn: func(_ context.Context, productID string, qty int) error {
			if productID != "product-1" || qty != 2 {
				t.Fatalf("reserve stock args = %q %d, want product-1 2", productID, qty)
			}
			return nil
		},
	}

	got, err := NewOrderService(txManager, orderRepo, itemRepo, productSvc, inventorySvc).CreateOrder(context.Background(), &dto.CreateOrderReq{
		CustomerID:      "customer-1",
		ShippingAddress: "Jakarta",
		Items:           []dto.OrderItemReq{{ProductID: "product-1", Quantity: 2}},
	})
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if !txManager.called {
		t.Fatal("CreateOrder must run writes inside a transaction")
	}
	if got.ID != "order-1" || got.OrderNumber != "ORD-20260720-0001" || got.TotalAmount != 20 {
		t.Fatalf("order response = %#v", got)
	}
}
