package usecase

type EmailUseCase interface {
	SendEmailVerification(recipientName, recipientEmail, activationCode string) error
}
