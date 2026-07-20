package service

import (
	"Dzaakk/simple-commerce/internal/order/dto"
	"Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/package/constant"
	dbtx "Dzaakk/simple-commerce/package/db/transactor"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"net/http"
	"time"
)

type OrderServiceImpl struct {
	transactor    dbtx.Transactor
	orderRepo     OrderRepository
	orderItemRepo OrderItemRepository
	productSvc    ProductService
	inventorySvc  InventoryService
}

func NewOrderService(transactor dbtx.Transactor, orderRepo OrderRepository, orderItemRepo OrderItemRepository, productSvc ProductService, inventorySvc InventoryService) *OrderServiceImpl {
	return &OrderServiceImpl{
		transactor:    transactor,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		productSvc:    productSvc,
		inventorySvc:  inventorySvc,
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

	orderNumber, err := s.orderRepo.GenerateOrderNumber(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	var (
		orderID string
		total   float64
		items   = make([]*model.OrderItem, 0, len(req.Items))
		status  = constant.OrderPending
		order   *model.Order
	)

	if err := s.transactor.WithinTx(ctx, func(txCtx context.Context) error {
		for _, item := range req.Items {
			if item.ProductID == "" || item.Quantity <= 0 {
				return response.NewAppError(http.StatusBadRequest, "invalid parameter item")
			}

			product, err := s.productSvc.FindByID(txCtx, item.ProductID)
			if err != nil {
				return err
			}
			if product == nil || !product.IsActive {
				return response.NewAppError(http.StatusNotFound, "product not found")
			}

			if err := s.inventorySvc.ReserveStock(txCtx, item.ProductID, item.Quantity); err != nil {
				return err
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

		order = &model.Order{
			OrderNumber:     orderNumber,
			CustomerID:      req.CustomerID,
			Status:          string(status),
			TotalAmount:     total,
			ShippingAddress: req.ShippingAddress,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		var err error
		orderID, err = s.orderRepo.Create(txCtx, order)
		if err != nil {
			return err
		}
		order.ID = orderID

		for _, item := range items {
			item.OrderID = orderID
		}

		return s.orderItemRepo.CreateBatch(txCtx, items)
	}); err != nil {
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

	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, response.NewAppError(http.StatusNotFound, "order not found")
	}
	if order.CustomerID != customerID {
		return nil, response.NewAppError(http.StatusUnauthorized, "unauthorized")
	}

	items, err := s.orderItemRepo.FindByOrderID(ctx, orderID)
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

	orders, err := s.orderRepo.FindByCustomerID(ctx, customerID, filter)
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

	return s.transactor.WithinTx(ctx, func(txCtx context.Context) error {
		order, err := s.orderRepo.FindByID(txCtx, orderID)
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

		items, err := s.orderItemRepo.FindByOrderID(txCtx, orderID)
		if err != nil {
			return err
		}

		for _, item := range items {
			if item == nil {
				continue
			}
			if err := s.inventorySvc.ReleaseStock(txCtx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}

		return s.orderRepo.UpdateStatus(txCtx, orderID, constant.OrderCancelled)
	})
}

func (s *OrderServiceImpl) UpdateOrderStatus(ctx context.Context, orderID string, status constant.OrderStatus) error {
	if orderID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter order id")
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, status)
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
