package route

import (
	"Dzaakk/simple-commerce/internal/auth/handler"
	"Dzaakk/simple-commerce/internal/auth/repository"
	"Dzaakk/simple-commerce/internal/auth/usecase"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repository"
	customerUsecase "Dzaakk/simple-commerce/internal/customer/usecase"
	emailUsecase "Dzaakk/simple-commerce/internal/email/usecase"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	sellerRepo "Dzaakk/simple-commerce/internal/seller/repository"
	shoppingCartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AuthRoutes struct {
	Handler            *handler.AuthHandler
	CustomerMiddleware *middleware.JWTCustomerMiddleware
	SellerMiddleware   *middleware.JWTSellerMiddleware
}

func NewAuthRoutes(handler *handler.AuthHandler, customerMiddleware *middleware.JWTCustomerMiddleware, sellerMiddleware *middleware.JWTSellerMiddleware) *AuthRoutes {
	return &AuthRoutes{
		Handler:            handler,
		CustomerMiddleware: customerMiddleware,
		SellerMiddleware:   sellerMiddleware,
	}
}

func (ar *AuthRoutes) Route(r *gin.RouterGroup) {
	apiGroup := r.Group("/api/v1")

	customer := apiGroup.Group("/customer")
	{
		customer.POST("/register", ar.Handler.RegistrationCustomer)
		customer.POST("/activate", ar.Handler.ActivationCustomer)
		customer.POST("/login", ar.Handler.LoginCustomer)
		customer.POST("/logout", ar.CustomerMiddleware.ValidateToken(), ar.Handler.Logout)
	}

	seller := apiGroup.Group("/seller")
	{
		seller.POST("/register", ar.Handler.RegistrationSeller)
		seller.POST("/activate", ar.Handler.ActivationSeller)
		seller.POST("/login", ar.Handler.LoginSeller)
		seller.POST("/logout", ar.SellerMiddleware.ValidateToken(), ar.Handler.Logout)
	}
}

func InitializedService(db *sql.DB, redis *redis.Client) *AuthRoutes {
	customerRepository := customerRepo.NewCustomerRepository(db)
	customerUsecase := customerUsecase.NewCustomerUseCase(customerRepository)

	sellerRepository := sellerRepo.NewSellerRepository(db)
	shoppingCartRepository := shoppingCartRepo.NewShoppingCartRepository(db)

	authCacheCustomer := repository.NewAuthCacheCustomerRepository(redis)
	authCacheSeller := repository.NewAuthCacheSellerRepository(redis)

	authUsecase := usecase.NewAuthUsecase(authCacheCustomer, authCacheSeller, customerUsecase, sellerRepository, shoppingCartRepository)
	mailer := emailUsecase.NewEmailUseCase()
	handler := handler.NewAtuhHandler(authUsecase, mailer)

	authCustomerToken := usecase.NewAuthCustomerTokenUsecase(authCacheCustomer)
	authSellerToken := usecase.NewAuthSellerTokenUsecase(authCacheSeller)

	customerMiddleware := middleware.NewJWTCustomerMiddleware(authCustomerToken)
	sellerMiddleware := middleware.NewJWTSellerMiddleware(authSellerToken)

	return NewAuthRoutes(handler, customerMiddleware, sellerMiddleware)
}
