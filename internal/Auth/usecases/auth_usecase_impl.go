package usecases

import (
	customerRepo "Dzaakk/simple-commerce/internal/customer/repositories"
)

type AuthUseCaseImpl struct {
	customerRepo customerRepo.CustomerRepository
}

func NewAuthUseCase(customerRepo customerRepo.CustomerRepository) AuthUseCase {
	return &AuthUseCaseImpl{customerRepo}
}

// CustomerLogin implements AuthUseCase.
func (a *AuthUseCaseImpl) CustomerLogin() {
	panic("unimplemented")
}

// CustomerRegistration implements AuthUseCase.
func (a *AuthUseCaseImpl) CustomerRegistration() {
	panic("unimplemented")
}
