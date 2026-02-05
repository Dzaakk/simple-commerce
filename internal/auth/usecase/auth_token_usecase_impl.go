package usecase

import "context"

type AuthCustomerTokenUsecaseImpl struct {
	AuthCacheCustomer AuthCacheCustomer
}

func NewAuthCustomerTokenUsecase(authCacheCustomer AuthCacheCustomer) CustomerTokenUsecase {
	return &AuthCustomerTokenUsecaseImpl{AuthCacheCustomer: authCacheCustomer}
}

func (a *AuthCustomerTokenUsecaseImpl) GetToken(ctx context.Context, email string) (*string, error) {
	return a.AuthCacheCustomer.GetToken(ctx, email)
}

type AuthSellerTokenUsecaseImpl struct {
	AuthCacheSeller AuthCacheSeller
}

func NewAuthSellerTokenUsecase(authCacheSeller AuthCacheSeller) SellerTokenUsecase {
	return &AuthSellerTokenUsecaseImpl{AuthCacheSeller: authCacheSeller}
}

func (a *AuthSellerTokenUsecaseImpl) GetToken(ctx context.Context, email string) (*string, error) {
	return a.AuthCacheSeller.GetToken(ctx, email)
}
