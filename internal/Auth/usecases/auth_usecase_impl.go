package usecases

import (
	model "Dzaakk/simple-commerce/internal/auth/models"
	repo "Dzaakk/simple-commerce/internal/auth/repositories"
	customerModel "Dzaakk/simple-commerce/internal/customer/models"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repositories"
	sellerModel "Dzaakk/simple-commerce/internal/seller/models"
	sellerRepo "Dzaakk/simple-commerce/internal/seller/repositories"
	shoppingCartModel "Dzaakk/simple-commerce/internal/shopping_cart/models"
	shoppingCartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repositories"
	template "Dzaakk/simple-commerce/package/templates"
	"context"
	"time"
)

type AuthUseCaseImpl struct {
	Repo             repo.AuthRepository
	CustomerRepo     customerRepo.CustomerRepository
	SellerRepo       sellerRepo.SellerRepository
	ShoppingCartRepo shoppingCartRepo.ShoppingCartRepository
}

func NewAuthUseCase(repo repo.AuthRepository, customerRepo customerRepo.CustomerRepository, sellerRepo sellerRepo.SellerRepository, shoppingCartRepo shoppingCartRepo.ShoppingCartRepository) AuthUseCase {
	return &AuthUseCaseImpl{repo, customerRepo, sellerRepo, shoppingCartRepo}
}

func (a *AuthUseCaseImpl) CustomerRegistration(ctx context.Context, data model.CustomerRegistration) (*int64, error) {

	hashedPassword, err := template.HashPassword(data.Password)
	if err != nil {
		return nil, err
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

	customerID, err := a.CustomerRepo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	NewShoppingCart := shoppingCartModel.TShoppingCart{
		CustomerID: int(customerID),
		Status:     "A",
		Base: template.Base{
			Created:   customer.Created,
			CreatedBy: "System",
		},
	}

	_, err = a.ShoppingCartRepo.Create(ctx, NewShoppingCart)
	if err != nil {
		return nil, err
	}

	codeActivation := GenerateActivationCode()
	newActivationCode := model.TCustomerActivationCode{
		CustomerID:     customerID,
		CodeActivation: codeActivation,
		IsUsed:         false,
		CreatedAt:      time.Now(),
	}

	err = a.Repo.InsertCustomerCodeActivation(ctx, newActivationCode)
	if err != nil {
		return nil, err
	}

	//send email

	return &customerID, nil
}

func (a *AuthUseCaseImpl) SellerRegistration(ctx context.Context, data model.SellerRegistration) (*int64, error) {
	hashedPassword, err := template.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	seller := sellerModel.TSeller{
		Username: data.Username,
		Email:    data.Email,
		Password: string(hashedPassword),
		Balance:  float64(0),
		Status:   "A",
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}

	sellerID, err := a.SellerRepo.Create(ctx, seller)
	if err != nil {
		return nil, err
	}

	codeActivation := GenerateActivationCode()
	newActivationCode := model.TSellerActivationCode{
		SellerID:       sellerID,
		CodeActivation: codeActivation,
		IsUsed:         false,
		CreatedAt:      time.Now(),
	}

	err = a.Repo.InsertSellerCodeActivation(ctx, newActivationCode)
	if err != nil {
		return nil, err
	}

	//send email

	return &sellerID, nil
}
