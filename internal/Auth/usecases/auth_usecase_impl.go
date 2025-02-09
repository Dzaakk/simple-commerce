package usecases

import (
	model "Dzaakk/simple-commerce/internal/auth/models"
	customerModel "Dzaakk/simple-commerce/internal/customer/models"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repositories"
	template "Dzaakk/simple-commerce/package/templates"
	"context"
	"time"
)

type AuthUseCaseImpl struct {
	customerRepo customerRepo.CustomerRepository
}

func NewAuthUseCase(customerRepo customerRepo.CustomerRepository) AuthUseCase {
	return &AuthUseCaseImpl{customerRepo}
}

func (a *AuthUseCaseImpl) CustomerRegistration(ctx context.Context, data model.CustomerRegistration) (int64, error) {

	hashedPassword, err := template.HashPassword(data.Password)
	if err != nil {
		return 0, err
	}

	customer := customerModel.TCustomers{
		Username:    data.Username,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    string(hashedPassword),
		Balance:     float64(10000000),
		Status:      "A",
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}

	customerId, err := a.customerRepo.Create(ctx, customer)
	if err != nil {
		return 0, err
	}
	return customerId, nil
}

func (a *AuthUseCaseImpl) CustomerLogin() {
	panic("unimplemented")
}
