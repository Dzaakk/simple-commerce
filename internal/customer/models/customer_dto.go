package models

type CustomerReq struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type CustomerRes struct {
	Id          string `json:"id"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Balance     string `json:"balance,omitempty"`
}

type BalanceUpdateReq struct {
	Id         string `json:"id" validate:"required,numeric,min=1"`
	ActionType string `json:"actionType" validate:"required"`
	Balance    string `json:"balance" validate:"required"`
}
type BalanceUpdateRes struct {
	BalanceOld CustomerBalance `json:"oldData"`
	BalanceNew CustomerBalance `json:"newData"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
