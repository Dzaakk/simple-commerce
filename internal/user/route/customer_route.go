package route

import (
	authRepo "Dzaakk/simple-commerce/internal/auth/repository"
	authUsecase "Dzaakk/simple-commerce/internal/auth/usecase"
	"Dzaakk/simple-commerce/internal/customer/handler"
	"Dzaakk/simple-commerce/internal/customer/repository"
	"Dzaakk/simple-commerce/internal/customer/usecase"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type CustomerRoutes struct {
	Handler       *handler.CustomerHandler
	JWTMiddleware *middleware.JWTCustomerMiddleware
}

func NewCustomerRoutes(handler *handler.CustomerHandler, jwtMiddleware *middleware.JWTCustomerMiddleware) *CustomerRoutes {
	return &CustomerRoutes{
		Handler:       handler,
		JWTMiddleware: jwtMiddleware,
	}
}

func (cr *CustomerRoutes) Route(r *gin.RouterGroup) {
	customerHandler := r.Group("/api/v1/customer")
	customerHandler.Use(cr.JWTMiddleware.ValidateToken())
	{
		customerHandler.GET("/find-all", cr.Handler.FindCustomerByID)
	}
}

func InitializedService(db *sql.DB, redis *redis.Client) *CustomerRoutes {
	repo := repository.NewCustomerRepository(db)
	usecase := usecase.NewCustomerUseCase(repo)
	handler := handler.NewCustomerHandler(usecase)

	authCachecustomer := authRepo.NewAuthCacheCustomerRepository(redis)
	authCustomerToken := authUsecase.NewAuthCustomerTokenUsecase(authCachecustomer)
	jwtMiddleware := middleware.NewJWTCustomerMiddleware(authCustomerToken)

	return NewCustomerRoutes(handler, jwtMiddleware)
}
