package service

import (
	"context"
	"testing"

	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
)

type customerRepositoryStub struct {
	create func(*model.Customer) (string, error)
	byID   func(string) (*model.Customer, error)
}

func (s customerRepositoryStub) Create(_ context.Context, customer *model.Customer) (string, error) {
	return s.create(customer)
}
func (s customerRepositoryStub) FindByEmail(context.Context, string) (*model.Customer, error) {
	return nil, nil
}
func (s customerRepositoryStub) FindByID(_ context.Context, id string) (*model.Customer, error) {
	return s.byID(id)
}

func TestCreateCustomerIsActiveImmediately(t *testing.T) {
	repo := customerRepositoryStub{create: func(customer *model.Customer) (string, error) {
		if customer.Status != string(constant.StatusActive) {
			t.Fatalf("status = %q", customer.Status)
		}
		if customer.Email != "customer@example.com" || customer.FullName != "Customer" {
			t.Fatalf("customer = %#v", customer)
		}
		return "customer-1", nil
	}}
	id, err := NewCustomerService(repo).Create(context.Background(), &dto.RegisterCustomerRequest{
		Email: " Customer@Example.COM ", FullName: " Customer ", PasswordHash: "hash",
	})
	if err != nil || id != "customer-1" {
		t.Fatalf("id = %q, error = %v", id, err)
	}
}

func TestFindByIDDoesNotExposePasswordHash(t *testing.T) {
	repo := customerRepositoryStub{
		create: func(*model.Customer) (string, error) { return "", nil },
		byID: func(string) (*model.Customer, error) {
			return &model.Customer{ID: "customer-1", Email: "customer@example.com", PasswordHash: "secret", Status: string(constant.StatusActive)}, nil
		},
	}
	result, err := NewCustomerService(repo).FindByID(context.Background(), "customer-1")
	if err != nil || result == nil || result.ID != "customer-1" {
		t.Fatalf("result = %#v, error = %v", result, err)
	}
}
