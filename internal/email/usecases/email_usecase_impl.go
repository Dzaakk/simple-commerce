package usecases

type EmailUsecCaseImpl struct {
}

func NewEmailUseCase() EmailUseCase {
	return &EmailUsecCaseImpl{}
}

// SendEmailVerification implements EmailUseCase.
func (e *EmailUsecCaseImpl) SendEmailVerification(recipientName, recipientEmail, activationCode string) error {
	panic("unimplemented")
}
