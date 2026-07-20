package service

import (
	orderModel "Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/internal/transaction/dto"
	txModel "Dzaakk/simple-commerce/internal/transaction/model"
	"Dzaakk/simple-commerce/package/constant"
	dbtx "Dzaakk/simple-commerce/package/db/transactor"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"net/http"
	"time"
)

type TransactionServiceImpl struct {
	transactor       dbtx.Transactor
	txRepo           TransactionRepository
	orderRepo        OrderRepository
	orderItemRepo    OrderItemRepository
	orderService     OrderService
	inventoryService InventoryService
}

func NewTransactionService(transactor dbtx.Transactor, txRepo TransactionRepository, orderRepo OrderRepository, orderItemRepo OrderItemRepository, orderService OrderService, inventoryService InventoryService) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		transactor:       transactor,
		txRepo:           txRepo,
		orderRepo:        orderRepo,
		orderItemRepo:    orderItemRepo,
		orderService:     orderService,
		inventoryService: inventoryService,
	}
}

func (s *TransactionServiceImpl) CreateTransaction(ctx context.Context, req *dto.CreateTransactionReq) (*dto.TransactionRes, error) {
	if req == nil {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid request")
	}
	if req.CustomerID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter customer id")
	}
	if req.OrderID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter order id")
	}
	if req.PaymentMethod == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter payment method")
	}

	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, response.NewAppError(http.StatusNotFound, "order not found")
	}
	if order.CustomerID != req.CustomerID {
		return nil, response.NewAppError(http.StatusUnauthorized, "unauthorized")
	}
	if order.Status != string(constant.OrderPending) {
		return nil, response.NewAppError(http.StatusConflict, "order status is not pending")
	}

	existing, err := s.txRepo.FindByOrderID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, response.NewAppError(http.StatusConflict, "transaction already exists")
	}

	txNumber, err := s.txRepo.GenerateTransactionNumber(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	data := &txModel.Transaction{
		OrderID:           order.ID,
		TransactionNumber: txNumber,
		PaymentMethod:     req.PaymentMethod,
		Status:            string(constant.TransactionPending),
		Amount:            order.TotalAmount,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	var id string
	if err := s.transactor.WithinTx(ctx, func(txCtx context.Context) error {
		var err error
		id, err = s.txRepo.Create(txCtx, data)
		return err
	}); err != nil {
		return nil, err
	}
	data.ID = id

	return toTransactionRes(data), nil
}

func (s *TransactionServiceImpl) GetTransactionByID(ctx context.Context, customerID, transactionID string) (*dto.TransactionRes, error) {
	if customerID == "" || transactionID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter")
	}

	txData, err := s.txRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	if txData == nil {
		return nil, response.NewAppError(http.StatusNotFound, "transaction not found")
	}

	order, err := s.orderRepo.FindByID(ctx, txData.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, response.NewAppError(http.StatusNotFound, "order not found")
	}
	if order.CustomerID != customerID {
		return nil, response.NewAppError(http.StatusUnauthorized, "unauthorized")
	}

	return toTransactionRes(txData), nil
}

func (s *TransactionServiceImpl) GetTransactionByOrderID(ctx context.Context, customerID, orderID string) (*dto.TransactionRes, error) {
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

	txData, err := s.txRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if txData == nil {
		return nil, response.NewAppError(http.StatusNotFound, "transaction not found")
	}

	return toTransactionRes(txData), nil
}

func (s *TransactionServiceImpl) HandlePaymentCallback(ctx context.Context, req *dto.PaymentCallbackReq) error {
	if req == nil {
		return response.NewAppError(http.StatusBadRequest, "invalid request")
	}
	if req.TransactionNumber == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter transaction number")
	}
	if req.Signature == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid signature")
	}

	if !verifySignature(req) {
		return response.NewAppError(http.StatusBadRequest, "invalid signature")
	}

	return s.transactor.WithinTx(ctx, func(txCtx context.Context) error {
		txData, err := s.txRepo.FindByTransactionNumber(txCtx, req.TransactionNumber)
		if err != nil {
			return err
		}
		if txData == nil {
			return response.NewAppError(http.StatusNotFound, "transaction not found")
		}

		currentStatus := constant.TransactionStatus(txData.Status)
		if isFinalTransactionStatus(currentStatus) {
			return nil
		}

		newStatus := req.Status
		if newStatus == "" {
			return response.NewAppError(http.StatusBadRequest, "invalid parameter status")
		}

		order, err := s.orderRepo.FindByID(txCtx, txData.OrderID)
		if err != nil {
			return err
		}
		if order == nil {
			return response.NewAppError(http.StatusNotFound, "order not found")
		}

		items, err := s.orderItemRepo.FindByOrderID(txCtx, order.ID)
		if err != nil {
			return err
		}

		paidAt := req.PaidAt
		if newStatus == constant.TransactionSuccess && paidAt == nil {
			now := time.Now()
			paidAt = &now
		}

		if err := s.txRepo.UpdateStatus(txCtx, txData.ID, newStatus, paidAt); err != nil {
			return err
		}

		switch newStatus {
		case constant.TransactionSuccess:
			if err := s.orderService.UpdateOrderStatus(txCtx, order.ID, constant.OrderConfirmed); err != nil {
				return err
			}
		case constant.TransactionFailed, constant.TransactionExpired:
			if err := s.orderService.UpdateOrderStatus(txCtx, order.ID, constant.OrderCancelled); err != nil {
				return err
			}
			if err := releaseOrderItems(txCtx, items, s.inventoryService); err != nil {
				return err
			}
		case constant.TransactionRefunded:
			if err := s.orderService.UpdateOrderStatus(txCtx, order.ID, constant.OrderCancelled); err != nil {
				return err
			}
			if err := releaseOrderItems(txCtx, items, s.inventoryService); err != nil {
				return err
			}
		default:
			// pending/processing: no-op for order and inventory
		}

		return nil
	})
}

func (s *TransactionServiceImpl) ExpireTransaction(ctx context.Context, transactionID string) error {
	if transactionID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter transaction id")
	}

	return s.transactor.WithinTx(ctx, func(txCtx context.Context) error {
		txData, err := s.txRepo.FindByID(txCtx, transactionID)
		if err != nil {
			return err
		}
		if txData == nil {
			return response.NewAppError(http.StatusNotFound, "transaction not found")
		}

		currentStatus := constant.TransactionStatus(txData.Status)
		if isFinalTransactionStatus(currentStatus) {
			return nil
		}

		order, err := s.orderRepo.FindByID(txCtx, txData.OrderID)
		if err != nil {
			return err
		}
		if order == nil {
			return response.NewAppError(http.StatusNotFound, "order not found")
		}

		items, err := s.orderItemRepo.FindByOrderID(txCtx, order.ID)
		if err != nil {
			return err
		}

		if err := s.txRepo.UpdateStatus(txCtx, txData.ID, constant.TransactionExpired, nil); err != nil {
			return err
		}
		if err := s.orderService.UpdateOrderStatus(txCtx, order.ID, constant.OrderCancelled); err != nil {
			return err
		}
		return releaseOrderItems(txCtx, items, s.inventoryService)
	})
}

func releaseOrderItems(ctx context.Context, items []*orderModel.OrderItem, inventorySvc InventoryService) error {
	for _, item := range items {
		if item == nil {
			continue
		}
		if err := inventorySvc.ReleaseStock(ctx, item.ProductID, item.Quantity); err != nil {
			return err
		}
	}
	return nil
}

func isFinalTransactionStatus(status constant.TransactionStatus) bool {
	switch status {
	case constant.TransactionSuccess, constant.TransactionFailed, constant.TransactionExpired, constant.TransactionRefunded:
		return true
	default:
		return false
	}
}

func verifySignature(req *dto.PaymentCallbackReq) bool {
	// TODO: implement real signature verification based on gateway spec.
	return req.Signature != ""
}

func toTransactionRes(data *txModel.Transaction) *dto.TransactionRes {
	if data == nil {
		return nil
	}

	status := constant.TransactionStatus(data.Status)
	return &dto.TransactionRes{
		ID:                data.ID,
		OrderID:           data.OrderID,
		TransactionNumber: data.TransactionNumber,
		PaymentMethod:     data.PaymentMethod,
		Status:            status,
		Amount:            data.Amount,
		PaidAt:            data.PaidAt,
		CreatedAt:         data.CreatedAt,
	}
}
