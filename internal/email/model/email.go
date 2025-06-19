package model

type BaseEmailReq struct {
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
type ActivationEmailReq struct {
	Email          string
	Username       string
	ActivationCode string
}
