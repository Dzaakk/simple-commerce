package usecases

import (
	model "Dzaakk/simple-commerce/internal/Auth/models"
	"context"
)

type AuthUseCase interface {
	CustomerRegistration(ctx context.Context, data model.CustomerRegistration) (int64, error)
	CustomerLogin()
}
