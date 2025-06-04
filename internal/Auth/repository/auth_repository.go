package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"context"
)

type AuthRepository interface {
	InsertCustomerCodeActivation(c context.Context, data model.TCustomerActivationCode) error
	FindCodeByCustomerID(c context.Context, customerID int64) (*model.TCustomerActivationCode, error)
	InsertSellerCodeActivation(c context.Context, data model.TSellerActivationCode) error
	FindCodeBySellerID(c context.Context, sellerID int64) (*model.TSellerActivationCode, error)
}
