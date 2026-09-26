package route

import (
	"database/sql"

	"Dzaakk/simple-commerce/internal/middleware"
	"Dzaakk/simple-commerce/internal/user/handler"
	"Dzaakk/simple-commerce/internal/user/repository"
	"Dzaakk/simple-commerce/internal/user/service"

	"github.com/gin-gonic/gin"
)

func Route(router *gin.RouterGroup, db *sql.DB) {
	customers := service.NewCustomerService(repository.NewCustomerRepository(db))
	h := handler.NewUserHandler(customers)
	router.GET("/api/v1/customer/me", middleware.Authenticate(), h.Me)
}
