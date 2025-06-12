package usecase

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	repo "Dzaakk/simple-commerce/internal/auth/repository"
	customerModel "Dzaakk/simple-commerce/internal/customer/model"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repository"
	sellerModel "Dzaakk/simple-commerce/internal/seller/model"
	sellerRepo "Dzaakk/simple-commerce/internal/seller/repository"
	shoppingCartModel "Dzaakk/simple-commerce/internal/shopping_cart/model"
	shoppingCartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/package/template"
	"Dzaakk/simple-commerce/package/util"
	"context"
	"database/sql"
	"errors"
	"os"
	"time"
)

type AuthUseCaseImpl struct {
	Cache            repo.AuthCacheRepository
	CustomerRepo     customerRepo.CustomerRepository
	SellerRepo       sellerRepo.SellerRepository
	ShoppingCartRepo shoppingCartRepo.ShoppingCartRepository
}

func NewAuthUseCase(cache repo.AuthCacheRepository, customerRepo customerRepo.CustomerRepository, sellerRepo sellerRepo.SellerRepository, shoppingCartRepo shoppingCartRepo.ShoppingCartRepository) AuthUseCase {
	return &AuthUseCaseImpl{cache, customerRepo, sellerRepo, shoppingCartRepo}
}

func (a *AuthUseCaseImpl) CustomerRegistration(ctx context.Context, data model.CustomerRegistrationReq) error {

	codeActivation := GenerateActivationCode()
	err := a.Cache.SetActivationCustomer(ctx, data.Email, codeActivation)
	if err != nil {
		return err
	}

	err = a.Cache.SetCustomerRegistration(ctx, data)
	if err != nil {
		return err
	}

	//send email

	return nil
}

func (a *AuthUseCaseImpl) CustomerActivation(ctx context.Context, req model.CustomerActivationReq) error {
	code, err := a.Cache.GetActivationCustomer(ctx, req.Email)
	if err != nil {
		return err
	}

	if code != req.ActivationCode {
		return errors.New("invalid activation code")
	}

	data, err := a.Cache.GetCustomerRegistration(ctx, req.Email)
	if err != nil {
		return err
	}

	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return err
	}

	customer := customerModel.TCustomers{
		Username:       data.Username,
		Email:          data.Email,
		PhoneNumber:    data.PhoneNumber,
		Password:       string(hashedPassword),
		Balance:        float64(10000000),
		Status:         1,
		Gender:         data.Gender,
		DateOfBirth:    data.DateOfBirth,
		LastLogin:      sql.NullTime{Time: time.Now(), Valid: true},
		ProfilePicture: "",
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}
	customerID, err := a.CustomerRepo.Create(ctx, customer)
	if err != nil {
		return err
	}

	NewShoppingCart := shoppingCartModel.TShoppingCart{
		CustomerID: int(customerID),
		Status:     template.StatusActive,
		Base: template.Base{
			Created:   customer.Created,
			CreatedBy: "System",
		},
	}

	_, err = a.ShoppingCartRepo.Create(ctx, NewShoppingCart)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthUseCaseImpl) CustomerLogin(ctx context.Context, req model.LoginReq) error {

	customer, err := a.CustomerRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if !util.CheckPasswordHash(req.Password, customer.Password) {
		return errors.New("invalid email or password")
	}

	tokenData := model.CustomerToken{
		ID:       customer.ID,
		Username: customer.Username,
		Email:    customer.Email,
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	jwtToken, err := util.GenerateJWTToken(secretKey, tokenData)
	if err != nil {
		return err
	}

	err = a.Cache.SetTokenCustomer(ctx, customer.Email, jwtToken)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthUseCaseImpl) SellerRegistration(ctx context.Context, data model.SellerRegistrationReq) (*int64, error) {
	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	seller := sellerModel.TSeller{
		Username: data.Username,
		Email:    data.Email,
		Password: string(hashedPassword),
		Balance:  float64(0),
		Status:   template.StatusActive,
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}

	sellerID, err := a.SellerRepo.Create(ctx, seller)
	if err != nil {
		return nil, err
	}

	// codeActivation := GenerateActivationCode()
	//send email

	return &sellerID, nil
}
