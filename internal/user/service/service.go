package service

import (
	"Dzaakk/simple-commerce/internal/user/domain"
	"Dzaakk/simple-commerce/internal/user/dto"
	"context"
)

type CustomerService interface {
	Create(ctx context.Context, req *dto.CreateReq) (string, error)
	Update(ctx context.Context, req *dto.UpdateReq) error
	FindByEmail(ctx context.Context, email string) (*domain.Customer, error)
	FindByID(ctx context.Context, customerID string) (*dto.CustomerRes, error)
}

type CustomerRepository interface {
	Create(ctx context.Context, data *domain.Customer) (string, error)
	Update(ctx context.Context, data *domain.Customer) (int64, error)
	FindByID(ctx context.Context, customerID string) (*domain.Customer, error)
	FindByEmail(ctx context.Context, email string) (*domain.Customer, error)
}

type SellerService interface {
	Create(ctx context.Context, req *dto.SellerCreateReq) (string, error)
	Update(ctx context.Context, req *dto.SellerUpdateReq) error
	FindByEmail(ctx context.Context, email string) (*domain.Seller, error)
	FindByID(ctx context.Context, sellerID string) (*dto.SellerRes, error)
}

type SellerRepository interface {
	Create(ctx context.Context, data *domain.Seller) (string, error)
	Update(ctx context.Context, data *domain.Seller) (int64, error)
	FindByID(ctx context.Context, sellerID string) (*domain.Seller, error)
	FindByEmail(ctx context.Context, email string) (*domain.Seller, error)
}
