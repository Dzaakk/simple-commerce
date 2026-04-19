package route

import (
	"Dzaakk/simple-commerce/internal/auth/handler"
	"Dzaakk/simple-commerce/internal/auth/repository"
	"Dzaakk/simple-commerce/internal/auth/service"
	emailQueue "Dzaakk/simple-commerce/internal/email/queue"
	emailService "Dzaakk/simple-commerce/internal/email/service"
	userrepo "Dzaakk/simple-commerce/internal/user/repository"
	userservice "Dzaakk/simple-commerce/internal/user/service"
	"Dzaakk/simple-commerce/package/rabbitmq"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AuthRoutes struct {
	Handler *handler.AuthHandler
}

func NewAuthRoutes(handler *handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		Handler: handler,
	}
}

func (ar *AuthRoutes) Route(r *gin.RouterGroup) {
	api := r.Group("/api/v1/auth")

	customer := api.Group("/customer")
	{
		customer.POST("/register", ar.Handler.RegisterCustomer)
	}

	seller := api.Group("/seller")
	{
		seller.POST("/register", ar.Handler.RegisterSeller)
	}

	api.GET("/verify-email", ar.Handler.VerifyEmail)
	api.GET("/login", ar.Handler.Login)
	api.POST("/refresh-token", ar.Handler.RefreshToken)
	api.POST("/logout", ar.Handler.Logout)
}

func InitializedService(db *sql.DB, redis *redis.Client, rabbit *rabbitmq.Client) *AuthRoutes {
	activationRepo := repository.NewActivationCodeRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)

	sellerRepo := userrepo.NewSellerRepository(db)
	customerRepo := userrepo.NewCustomerRepository(db)

	customerService := userservice.NewCustomerService(customerRepo)
	sellerService := userservice.NewSellerService(sellerRepo)
	emailService := emailService.NewEmailService()
	emailPublisher := emailQueue.NewRabbitPublisher(rabbit)

	service := service.NewAuthService(db, customerService, sellerService, emailService, emailPublisher, activationRepo, refreshRepo)

	hander := handler.NewAuthHandler(service)
	return NewAuthRoutes(hander)
}
