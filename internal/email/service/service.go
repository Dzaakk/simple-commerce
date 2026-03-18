package service

import (
	"Dzaakk/simple-commerce/internal/email/model"
	"context"
)

type EmailService interface {
	SendEmailVerification(ctx context.Context, req model.VerificationEmailReq) error
}
