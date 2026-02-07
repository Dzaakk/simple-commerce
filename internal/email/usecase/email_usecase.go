package usecase

import (
	"Dzaakk/simple-commerce/internal/email/model"
	"context"
)

type EmailUsecase interface {
	SendEmailActivation(ctx context.Context, req model.ActivationEmailReq) error
}
