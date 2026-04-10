package service

import (
	orderModel "Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/internal/transaction/dto"
	txModel "Dzaakk/simple-commerce/internal/transaction/model"
	"Dzaakk/simple-commerce/package/constant"
	"context"
	"database/sql"
	"errors"
	"time"
)

type TransactionServiceImpl struct {
	DB              *sql.DB
	TxRepo          TransactionRepository
	OrderRepo       OrderRepository
	OrderItemRepo   OrderItemRepository
	OrderService    OrderService
	InventoryService InventoryService
}

func NewTransactionService(db *sql.DB, txRepo TransactionRepository, orderRepo OrderRepository, orderItemRepo OrderItemRepository, orderService OrderService, inventoryService InventoryService) TransactionService {
	return &TransactionServiceImpl{
		DB:               db,
		TxRepo:           txRepo,
		OrderRepo:        orderRepo,
		OrderItemRepo:    orderItemRepo,
		OrderService:     orderService,
		InventoryService: inventoryService,
	}
}

func (s *TransactionServiceImpl) CreateTransaction(ctx context.Context, req *dto.CreateTransactionReq) (*dto.TransactionRes, error) {
	if req == nil {
		return nil, errors.New("invalid request")
	}
	if req.CustomerID == "" {
		return nil, errors.New("invalid parameter customer id")
	}
	if req.OrderID == "" {
		return nil, errors.New("invalid parameter order id")
	}
	if req.PaymentMethod == "" {
		return nil, errors.New("invalid parameter payment method")
	}

	order, err := s.OrderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.CustomerID != req.CustomerID {
		return nil, errors.New("unauthorized")
	}
	if order.Status != string(constant.OrderPending) {
		return nil, errors.New("order status is not pending")
	}

	existing, err := s.TxRepo.FindByOrderID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("transaction already exists")
	}

	txNumber, err := s.TxRepo.GenerateTransactionNumber(ctx)
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

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	id, err := s.TxRepo.Create(ctx, tx, data)
	if err != nil {
		return nil, err
	}
	data.ID = id

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return toTransactionRes(data), nil
}

func (s *TransactionServiceImpl) GetTransactionByID(ctx context.Context, customerID, transactionID string) (*dto.TransactionRes, error) {
	if customerID == "" || transactionID == "" {
		return nil, errors.New("invalid parameter")
	}

	txData, err := s.TxRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	if txData == nil {
		return nil, errors.New("transaction not found")
	}

	order, err := s.OrderRepo.FindByID(ctx, txData.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.CustomerID != customerID {
		return nil, errors.New("unauthorized")
	}

	return toTransactionRes(txData), nil
}

func (s *TransactionServiceImpl) GetTransactionByOrderID(ctx context.Context, customerID, orderID string) (*dto.TransactionRes, error) {
	if customerID == "" || orderID == "" {
		return nil, errors.New("invalid parameter")
	}

	order, err := s.OrderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.CustomerID != customerID {
		return nil, errors.New("unauthorized")
	}

	txData, err := s.TxRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if txData == nil {
		return nil, errors.New("transaction not found")
	}

	return toTransactionRes(txData), nil
}

func (s *TransactionServiceImpl) HandlePaymentCallback(ctx context.Context, req *dto.PaymentCallbackReq) error {
	if req == nil {
		return errors.New("invalid request")
	}
	if req.TransactionNumber == "" {
		return errors.New("invalid parameter transaction number")
	}
	if req.Signature == "" {
		return errors.New("invalid signature")
	}

	if !verifySignature(req) {
		return errors.New("invalid signature")
	}

	txData, err := s.TxRepo.FindByTransactionNumber(ctx, req.TransactionNumber)
	if err != nil {
		return err
	}
	if txData == nil {
		return errors.New("transaction not found")
	}

	currentStatus := constant.TransactionStatus(txData.Status)
	if isFinalTransactionStatus(currentStatus) {
		return nil
	}

	newStatus := req.Status
	if newStatus == "" {
		return errors.New("invalid parameter status")
	}

	order, err := s.OrderRepo.FindByID(ctx, txData.OrderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	items, err := s.OrderItemRepo.FindByOrderID(ctx, order.ID)
	if err != nil {
		return err
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	paidAt := req.PaidAt
	if newStatus == constant.TransactionSuccess && paidAt == nil {
		now := time.Now()
		paidAt = &now
	}

	if err := s.TxRepo.UpdateStatus(ctx, tx, txData.ID, newStatus, paidAt); err != nil {
		return err
	}

	switch newStatus {
	case constant.TransactionSuccess:
		if err := s.OrderService.UpdateOrderStatus(ctx, tx, order.ID, constant.OrderConfirmed); err != nil {
			return err
		}
	case constant.TransactionFailed, constant.TransactionExpired:
		if err := s.OrderService.UpdateOrderStatus(ctx, tx, order.ID, constant.OrderCancelled); err != nil {
			return err
		}
		if err := releaseOrderItems(ctx, tx, items, s.InventoryService); err != nil {
			return err
		}
	case constant.TransactionRefunded:
		if err := s.OrderService.UpdateOrderStatus(ctx, tx, order.ID, constant.OrderCancelled); err != nil {
			return err
		}
		if err := releaseOrderItems(ctx, tx, items, s.InventoryService); err != nil {
			return err
		}
	default:
		// pending/processing: no-op for order and inventory
	}

	return tx.Commit()
}

func (s *TransactionServiceImpl) ExpireTransaction(ctx context.Context, transactionID string) error {
	if transactionID == "" {
		return errors.New("invalid parameter transaction id")
	}

	txData, err := s.TxRepo.FindByID(ctx, transactionID)
	if err != nil {
		return err
	}
	if txData == nil {
		return errors.New("transaction not found")
	}

	currentStatus := constant.TransactionStatus(txData.Status)
	if isFinalTransactionStatus(currentStatus) {
		return nil
	}

	order, err := s.OrderRepo.FindByID(ctx, txData.OrderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	items, err := s.OrderItemRepo.FindByOrderID(ctx, order.ID)
	if err != nil {
		return err
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.TxRepo.UpdateStatus(ctx, tx, txData.ID, constant.TransactionExpired, nil); err != nil {
		return err
	}
	if err := s.OrderService.UpdateOrderStatus(ctx, tx, order.ID, constant.OrderCancelled); err != nil {
		return err
	}
	if err := releaseOrderItems(ctx, tx, items, s.InventoryService); err != nil {
		return err
	}

	return tx.Commit()
}

func releaseOrderItems(ctx context.Context, tx *sql.Tx, items []*orderModel.OrderItem, inventorySvc InventoryService) error {
	for _, item := range items {
		if item == nil {
			continue
		}
		if err := inventorySvc.ReleaseStock(ctx, tx, item.ProductID, item.Quantity); err != nil {
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
