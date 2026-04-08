package dto

import "Dzaakk/simple-commerce/package/constant"

type OrderFilter struct {
	Status *constant.OrderStatus
	Page   int
	Limit  int
}
