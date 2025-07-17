package handler

import "Dzaakk/simple-commerce/internal/order/usecase"

type OrderHandler struct {
	Usecase usecase.OrderUseCase
}

func NewOrderHandler(usecase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{Usecase: usecase}
}
