package usecase

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	customerModel "Dzaakk/simple-commerce/internal/customer/model"
	customerUsecase "Dzaakk/simple-commerce/internal/customer/usecase"
	eModel "Dzaakk/simple-commerce/internal/email/model"
	sellerModel "Dzaakk/simple-commerce/internal/seller/model"
	sellerRepo "Dzaakk/simple-commerce/internal/seller/repository"
	shoppingCartModel "Dzaakk/simple-commerce/internal/shopping_cart/model"
	shoppingCartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/package/template"
	"Dzaakk/simple-commerce/package/util"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecaseImpl struct {
	AuthCacheCustomer AuthCacheCustomer
	AuthCacheSeller   AuthCacheSeller
	CustomerUsecase   customerUsecase.CustomerUsecase
	SellerRepo        sellerRepo.SellerRepository
	ShoppingCartRepo  shoppingCartRepo.ShoppingCartRepository
}

func NewAuthUsecase(authCacheCustomer AuthCacheCustomer, authCacheSeller AuthCacheSeller, customerUsecase customerUsecase.CustomerUsecase, sellerRepo sellerRepo.SellerRepository, shoppingCartRepo shoppingCartRepo.ShoppingCartRepository) *AuthUsecaseImpl {
	return &AuthUsecaseImpl{
		AuthCacheCustomer: authCacheCustomer,
		AuthCacheSeller:   authCacheSeller,
		CustomerUsecase:   customerUsecase,
		SellerRepo:        sellerRepo,
		ShoppingCartRepo:  shoppingCartRepo,
	}
}

func (a *AuthUsecaseImpl) RegistrationCustomer(ctx context.Context, data model.CustomerRegistrationReq) (*eModel.ActivationEmailReq, error) {

	activationCode := generateActivationCode()
	err := a.AuthCacheCustomer.SetActivation(ctx, data.Email, activationCode)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	data.Password = string(hashedPassword)

	err = a.AuthCacheCustomer.SetRegistration(ctx, data)
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

func (a *AuthUsecaseImpl) ActivationCustomer(ctx context.Context, req model.ActivationReq) error {
	code, err := a.AuthCacheCustomer.GetActivation(ctx, req.Email)
	if err != nil {
		return err
	}

	if code != req.ActivationCode {
		return errors.New("invalid activation code")
	}

	data, err := a.AuthCacheCustomer.GetRegistration(ctx, req.Email)
	if err != nil {
		return err
	}

	customer := &customerModel.CreateReq{
		Username:    data.Username,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    data.Password,
		DateOfBirth: data.DateOfBirth,
		// Balance:        float64(10000000),
		// Status:         1,
		// Gender:         gender,
		// DateOfBirth:    sql.NullTime{Valid: true, Time: date},
		// LastLogin:      sql.NullTime{Time: time.Now(), Valid: true},
		// ProfilePicture: "",
		// Base: template.Base{
		// 	Created:   time.Now(),
		// 	CreatedBy: "system",
		// },
	}
	customerID, err := a.CustomerUsecase.Create(ctx, customer)
	if err != nil {
		return err
	}

	NewShoppingCart := shoppingCartModel.TShoppingCart{
		CustomerID: int(customerID),
		Status:     template.StatusActive,
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "System",
		},
	}

	_, err = a.ShoppingCartRepo.Create(ctx, NewShoppingCart)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthUsecaseImpl) LoginCustomer(ctx context.Context, req model.LoginReq) error {

	customer, err := a.CustomerUsecase.FindByEmail(ctx, req.Email)
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

	err = a.AuthCacheCustomer.SetToken(ctx, customer.Email, jwtToken)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthUsecaseImpl) RegistrationSeller(ctx context.Context, data model.SellerRegistrationReq) (*eModel.ActivationEmailReq, error) {
	activationCode := generateActivationCode()
	err := a.AuthCacheSeller.SetActivation(ctx, data.Email, activationCode)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	data.Password = string(hashedPassword)

	err = a.AuthCacheSeller.SetRegistration(ctx, data)
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
func (a *AuthUsecaseImpl) ActivationSeller(ctx context.Context, req model.ActivationReq) error {
	code, err := a.AuthCacheSeller.GetActivation(ctx, req.Email)
	if err != nil {
		return err
	}

	if code != req.ActivationCode {
		return errors.New("invalid activation code")
	}

	data, err := a.AuthCacheSeller.GetRegistration(ctx, req.Email)
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
func (a *AuthUsecaseImpl) LoginSeller(ctx context.Context, req model.LoginReq) error {
	seller, err := a.SellerRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if !util.CheckPasswordHash(req.Password, seller.Password) {
		return errors.New("invalid email or password")
	}

	tokenData := model.SellerToken{
		ID:       int(seller.ID),
		Username: seller.Username,
		Email:    seller.Email,
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

	err = a.AuthCacheSeller.SetToken(ctx, seller.Email, jwtToken)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthUsecaseImpl) Logout(ctx context.Context, email string, role string) error {

	switch role {
	case template.RoleCustomer:
		return a.AuthCacheCustomer.DeleteToken(ctx, email)
	case template.RoleSeller:
		return a.AuthCacheSeller.DeleteToken(ctx, email)
	}

	return nil
}
