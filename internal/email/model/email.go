package model

type EmailRequest struct {
	Sender      Sender      `json:"sender"`
	To          []Recipient `json:"to"`
	Subject     string      `json:"subject"`
	HTMLContent string      `json:"htmlContent"`
}

type Sender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type Recipient struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
