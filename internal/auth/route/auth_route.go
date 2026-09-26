package route

import (
	"database/sql"

	"Dzaakk/simple-commerce/internal/auth/handler"
	authrepo "Dzaakk/simple-commerce/internal/auth/repository"
	authservice "Dzaakk/simple-commerce/internal/auth/service"
	userrepo "Dzaakk/simple-commerce/internal/user/repository"
	userservice "Dzaakk/simple-commerce/internal/user/service"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	handler *handler.AuthHandler
}

func New(db *sql.DB) *AuthRoutes {
	customers := userservice.NewCustomerService(userrepo.NewCustomerRepository(db))
	tokens := authrepo.NewRefreshTokenRepository(db)
	return &AuthRoutes{handler: handler.NewAuthHandler(authservice.NewAuthService(customers, tokens))}
}

func (r *AuthRoutes) Route(router *gin.RouterGroup) {
	auth := router.Group("/api/v1/auth")
	auth.POST("/customer/register", r.handler.RegisterCustomer)
	auth.POST("/login", r.handler.Login)
	auth.POST("/refresh-token", r.handler.RefreshToken)
	auth.POST("/logout", r.handler.Logout)
}
