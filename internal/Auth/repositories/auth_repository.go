package repositories

import (
	model "Dzaakk/simple-commerce/internal/auth/models"
	"context"
)

type AuthRepository interface {
	CreateCode(ctx context.Context, data model.TActivationCode) error
	FindCodeByCustomerId(ctx context.Context, id int64) (*model.TActivationCode, error)
}
