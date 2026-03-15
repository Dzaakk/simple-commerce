package service

import (
	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	"context"
	"database/sql"
)

type CustomerService interface {
	Create(ctx context.Context, req *dto.RegisterCustomerRequest) (string, error)
	Update(ctx context.Context, req *dto.UpdateReq) error
	FindByEmail(ctx context.Context, email string) (*model.Customer, error)
	FindByID(ctx context.Context, customerID string) (*dto.CustomerRes, error)
	UpdateStatus(ctx context.Context, customerID string, status constant.UserStatus) error
	UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, customerID string, status constant.UserStatus) error
}

type CustomerRepository interface {
	Create(ctx context.Context, data *model.Customer) (string, error)
	Update(ctx context.Context, data *model.Customer) (int64, error)
	FindByID(ctx context.Context, customerID string) (*model.Customer, error)
	FindByEmail(ctx context.Context, email string) (*model.Customer, error)
	UpdateStatus(ctx context.Context, customerID string, status constant.UserStatus) error
	UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, customerID string, status constant.UserStatus) error
}

type SellerService interface {
	Create(ctx context.Context, req *dto.RegisterSellerRequest) (string, error)
	Update(ctx context.Context, req *dto.SellerUpdateReq) error
	FindByEmail(ctx context.Context, email string) (*model.Seller, error)
	FindByID(ctx context.Context, sellerID string) (*dto.SellerRes, error)
	UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, sellerID string, status constant.UserStatus) error
}

type SellerRepository interface {
	Create(ctx context.Context, data *model.Seller) (string, error)
	Update(ctx context.Context, data *model.Seller) (int64, error)
	FindByID(ctx context.Context, sellerID string) (*model.Seller, error)
	FindByEmail(ctx context.Context, email string) (*model.Seller, error)
	UpdateStatus(ctx context.Context, sellerID string, status constant.UserStatus) error
	UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, sellerID string, status constant.UserStatus) error
}
