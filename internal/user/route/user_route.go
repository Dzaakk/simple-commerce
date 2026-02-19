package route

import (
	authRepo "Dzaakk/simple-commerce/internal/auth/repository"
	authUsecase "Dzaakk/simple-commerce/internal/auth/usecase"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	"Dzaakk/simple-commerce/internal/user/handler"
	"Dzaakk/simple-commerce/internal/user/repository"
	"Dzaakk/simple-commerce/internal/user/service"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type UserRoutes struct {
	Handler       *handler.UserHandler
	JWTMiddleware *middleware.JWTCustomerMiddleware
}

func NewuserRoutes(handler *handler.UserHandler, jwtMiddleware *middleware.JWTCustomerMiddleware) *UserRoutes {
	return &UserRoutes{
		Handler:       handler,
		JWTMiddleware: jwtMiddleware,
	}
}

func (ur *UserRoutes) Route(r *gin.RouterGroup) {
	userHandler := r.Group("/api/v1/customer")
	userHandler.Use(ur.JWTMiddleware.ValidateToken())
	{
		userHandler.GET("/:id", ur.Handler.FindCustomerByID)
	}
}

func InitializedService(db *sql.DB, redis *redis.Client) *UserRoutes {
	repo := repository.NewCustomerRepository(db)
	service := service.NewCustomerService(repo)
	handler := handler.NewUserHandler(service)

	authCacheCustomer := authRepo.NewAuthCacheCustomerRepository(redis)
	authCustomerToken := authUsecase.NewAuthCustomerTokenUsecase(authCacheCustomer)
	jwtMiddleware := middleware.NewJWTCustomerMiddleware(authCustomerToken)

	return NewuserRoutes(handler, jwtMiddleware)
}
