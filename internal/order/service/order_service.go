package service

import (
	"Dzaakk/simple-commerce/internal/order/dto"
	"Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"net/http"
	"time"
)

type OrderServiceImpl struct {
	DB            *sql.DB
	OrderRepo     OrderRepository
	OrderItemRepo OrderItemRepository
	ProductSvc    ProductService
	InventorySvc  InventoryService
}

func NewOrderService(db *sql.DB, orderRepo OrderRepository, orderItemRepo OrderItemRepository, productSvc ProductService, inventorySvc InventoryService) OrderService {
	return &OrderServiceImpl{
		DB:            db,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		ProductSvc:    productSvc,
		InventorySvc:  inventorySvc,
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *dto.CreateOrderReq) (*dto.OrderRes, error) {
	if req == nil {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid request")
	}
	if req.CustomerID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter customer id")
	}
	if req.ShippingAddress == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter shipping address")
	}
	if len(req.Items) == 0 {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter items")
	}

	orderNumber, err := s.OrderRepo.GenerateOrderNumber(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()

	var (
		total  float64
		items  = make([]*model.OrderItem, 0, len(req.Items))
		status = constant.OrderPending
	)

	for _, item := range req.Items {
		if item.ProductID == "" || item.Quantity <= 0 {
			return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter item")
		}

		product, err := s.ProductSvc.FindByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil || !product.IsActive {
			return nil, response.NewAppError(http.StatusNotFound, "product not found")
		}

		if err := s.InventorySvc.ReserveStock(ctx, tx, item.ProductID, item.Quantity); err != nil {
			return nil, err
		}

		subtotal := product.Price * float64(item.Quantity)
		total += subtotal

		items = append(items, &model.OrderItem{
			ProductID: item.ProductID,
			SellerID:  product.SellerID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Subtotal:  subtotal,
			CreatedAt: now,
		})
	}

	order := &model.Order{
		OrderNumber:     orderNumber,
		CustomerID:      req.CustomerID,
		Status:          string(status),
		TotalAmount:     total,
		ShippingAddress: req.ShippingAddress,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	orderID, err := s.OrderRepo.Create(ctx, tx, order)
	if err != nil {
		return nil, err
	}
	order.ID = orderID

	for _, item := range items {
		item.OrderID = orderID
	}

	if err := s.OrderItemRepo.CreateBatch(ctx, tx, items); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.OrderRes{
		ID:              orderID,
		OrderNumber:     orderNumber,
		Status:          status,
		TotalAmount:     total,
		ShippingAddress: req.ShippingAddress,
		CreatedAt:       order.CreatedAt,
	}, nil
}

func (s *OrderServiceImpl) GetOrderByID(ctx context.Context, customerID string, orderID string) (*dto.OrderDetailRes, error) {
	if customerID == "" || orderID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter")
	}

	order, err := s.OrderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, response.NewAppError(http.StatusNotFound, "order not found")
	}
	if order.CustomerID != customerID {
		return nil, response.NewAppError(http.StatusUnauthorized, "unauthorized")
	}

	items, err := s.OrderItemRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	resItems := make([]dto.OrderItemRes, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		resItems = append(resItems, dto.OrderItemRes{
			ProductID: item.ProductID,
			SellerID:  item.SellerID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.Subtotal,
		})
	}

	return &dto.OrderDetailRes{
		OrderRes: toOrderRes(order),
		Items:    resItems,
	}, nil
}

func (s *OrderServiceImpl) GetOrdersByCustomer(ctx context.Context, customerID string, filter dto.OrderFilter) ([]*dto.OrderRes, error) {
	if customerID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter customer id")
	}

	orders, err := s.OrderRepo.FindByCustomerID(ctx, customerID, filter)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return []*dto.OrderRes{}, nil
	}

	result := make([]*dto.OrderRes, 0, len(orders))
	for _, order := range orders {
		if order == nil {
			continue
		}
		res := toOrderRes(order)
		result = append(result, &res)
	}

	return result, nil
}

func (s *OrderServiceImpl) CancelOrder(ctx context.Context, customerID string, orderID string) error {
	if customerID == "" || orderID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter")
	}

	order, err := s.OrderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return response.NewAppError(http.StatusNotFound, "order not found")
	}
	if order.CustomerID != customerID {
		return response.NewAppError(http.StatusUnauthorized, "unauthorized")
	}
	if order.Status != string(constant.OrderPending) {
		return response.NewAppError(http.StatusConflict, "order status is not pending")
	}

	items, err := s.OrderItemRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		if item == nil {
			continue
		}
		if err := s.InventorySvc.ReleaseStock(ctx, tx, item.ProductID, item.Quantity); err != nil {
			return err
		}
	}

	if err := s.OrderRepo.UpdateStatus(ctx, tx, orderID, constant.OrderCancelled); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *OrderServiceImpl) UpdateOrderStatus(ctx context.Context, tx *sql.Tx, orderID string, status constant.OrderStatus) error {
	if orderID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter order id")
	}
	if tx == nil {
		return response.NewAppError(http.StatusInternalServerError, "internal server error")
	}

	return s.OrderRepo.UpdateStatus(ctx, tx, orderID, status)
}

func toOrderRes(order *model.Order) dto.OrderRes {
	status := constant.OrderStatus(order.Status)
	return dto.OrderRes{
		ID:              order.ID,
		OrderNumber:     order.OrderNumber,
		Status:          status,
		TotalAmount:     order.TotalAmount,
		ShippingAddress: order.ShippingAddress,
		CreatedAt:       order.CreatedAt,
	}
}
