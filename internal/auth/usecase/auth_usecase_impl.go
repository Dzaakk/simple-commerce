package usecase

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	repo "Dzaakk/simple-commerce/internal/auth/repository"
	customerModel "Dzaakk/simple-commerce/internal/customer/model"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repository"
	eModel "Dzaakk/simple-commerce/internal/email/model"
	sellerModel "Dzaakk/simple-commerce/internal/seller/model"
	sellerRepo "Dzaakk/simple-commerce/internal/seller/repository"
	shoppingCartModel "Dzaakk/simple-commerce/internal/shopping_cart/model"
	shoppingCartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/package/template"
	"Dzaakk/simple-commerce/package/util"
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCaseImpl struct {
	CustomerCache    repo.AuthCacheCustomer
	SellerCache      repo.AuthCacheSeller
	CustomerRepo     customerRepo.CustomerRepository
	SellerRepo       sellerRepo.SellerRepository
	ShoppingCartRepo shoppingCartRepo.ShoppingCartRepository
}

func NewAuthUseCase(customerCache repo.AuthCacheCustomer, sellerCache repo.AuthCacheSeller, customerRepo customerRepo.CustomerRepository, sellerRepo sellerRepo.SellerRepository, shoppingCartRepo shoppingCartRepo.ShoppingCartRepository) AuthUseCase {
	return &AuthUseCaseImpl{customerCache, sellerCache, customerRepo, sellerRepo, shoppingCartRepo}
}

func (a *AuthUseCaseImpl) RegistrationCustomer(ctx context.Context, data model.CustomerRegistrationReq) (*eModel.ActivationEmailReq, error) {

	activationCode := GenerateActivationCode()
	err := a.CustomerCache.SetActivationCustomer(ctx, data.Email, activationCode)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	data.Password = string(hashedPassword)

	err = a.CustomerCache.SetCustomerRegistration(ctx, data)
	if err != nil {
		return nil, err
	}

	email := eModel.ActivationEmailReq{
		Email:          data.Email,
		Username:       data.Username,
		ActivationCode: activationCode,
	}

	return &email, nil
}

func (a *AuthUseCaseImpl) ActivationCustomer(ctx context.Context, req model.ActivationReq) error {
	code, err := a.CustomerCache.GetActivationCustomer(ctx, req.Email)
	if err != nil {
		return err
	}

	if code != req.ActivationCode {
		return errors.New("invalid activation code")
	}

	data, err := a.CustomerCache.GetCustomerRegistration(ctx, req.Email)
	if err != nil {
		return err
	}

	gender, err := strconv.Atoi(data.Gender)
	if err != nil {
		return err
	}

	date, err := time.Parse(data.DateOfBirth, template.FormatDate)
	if err != nil {
		return err
	}

	customer := customerModel.TCustomers{
		Username:       data.Username,
		Email:          data.Email,
		PhoneNumber:    data.PhoneNumber,
		Password:       data.Password,
		Balance:        float64(10000000),
		Status:         1,
		Gender:         gender,
		DateOfBirth:    sql.NullTime{Valid: true, Time: date},
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

func (a *AuthUseCaseImpl) LoginCustomer(ctx context.Context, req model.LoginReq) error {

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
		Role:     template.RoleCustomer,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken, err := util.GenerateToken(tokenData)
	if err != nil {
		return err
	}

	err = a.CustomerCache.SetTokenCustomer(ctx, customer.Email, jwtToken)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthUseCaseImpl) RegistrationSeller(ctx context.Context, data model.SellerRegistrationReq) (*eModel.ActivationEmailReq, error) {
	activationCode := GenerateActivationCode()
	err := a.SellerCache.SetActivationSeller(ctx, data.Email, activationCode)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	data.Password = string(hashedPassword)

	err = a.SellerCache.SetSellerRegistration(ctx, data)
	if err != nil {
		return nil, err
	}

	email := eModel.ActivationEmailReq{
		Email:          data.Email,
		Username:       data.Username,
		ActivationCode: activationCode,
	}

	return &email, nil
}
func (a *AuthUseCaseImpl) ActivationSeller(ctx context.Context, req model.ActivationReq) error {
	code, err := a.SellerCache.GetActivationSeller(ctx, req.Email)
	if err != nil {
		return err
	}

	if code != req.ActivationCode {
		return errors.New("invalid activation code")
	}

	data, err := a.SellerCache.GetSellerRegistration(ctx, req.Email)
	if err != nil {
		return err
	}

	seller := sellerModel.TSeller{
		Username:       data.Username,
		Email:          data.Email,
		PhoneNumber:    data.PhoneNumber,
		Password:       data.Password,
		Balance:        float64(10000000),
		Status:         1,
		StoreName:      data.StoreName,
		Address:        data.Address,
		ProfilePicture: "",
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}

	_, err = a.SellerRepo.Create(ctx, seller)
	if err != nil {
		return err
	}

	return nil
}
func (a *AuthUseCaseImpl) LoginSeller(ctx context.Context, req model.LoginReq) error {
	customer, err := a.CustomerRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if !util.CheckPasswordHash(req.Password, customer.Password) {
		return errors.New("invalid email or password")
	}

	tokenData := model.SellerToken{
		ID:       customer.ID,
		Username: customer.Username,
		Email:    customer.Email,
		Role:     template.RoleCustomer,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken, err := util.GenerateToken(tokenData)
	if err != nil {
		return err
	}

	err = a.SellerCache.SetTokenSeller(ctx, customer.Email, jwtToken)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthUseCaseImpl) Logout(ctx context.Context, email string, role string) error {

	switch role {
	case template.RoleCustomer:
		return a.CustomerCache.DeleteTokenCustomer(ctx, email)
	case template.RoleSeller:
		return a.SellerCache.DeleteTokenSeller(ctx, email)
	}

	return nil
}
