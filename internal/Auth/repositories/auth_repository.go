package repositories

import (
	model "Dzaakk/simple-commerce/internal/auth/models"
	"context"
)

type AuthRepository interface {
	InsertCustomerCodeActivation(c context.Context, data model.TCustomerActivationCode) error
	FindCodeByCustomerId(c context.Context, id int64) (*model.TCustomerActivationCode, error)
	InsertSellerCodeActivation(c context.Context, data model.TSellerActivationCode) error
	FindCodeBySellerId(c context.Context, id int64) (*model.TSellerActivationCode, error)
}
