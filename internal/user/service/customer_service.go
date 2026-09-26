package service

import (
	"context"
	"strings"

	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/model"
)

type CustomerServiceImpl struct {
	repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) *CustomerServiceImpl {
	return &CustomerServiceImpl{repo: repo}
}

func (s *CustomerServiceImpl) Create(ctx context.Context, req *dto.RegisterCustomerRequest) (string, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.FullName = strings.TrimSpace(req.FullName)
	return s.repo.Create(ctx, req.ToModel())
}

func (s *CustomerServiceImpl) FindByEmail(ctx context.Context, email string) (*model.Customer, error) {
	return s.repo.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
}

func (s *CustomerServiceImpl) FindByID(ctx context.Context, id string) (*dto.CustomerResponse, error) {
	customer, err := s.repo.FindByID(ctx, id)
	if err != nil || customer == nil {
		return nil, err
	}
	result := dto.ToCustomerResponse(customer)
	return &result, nil
}
