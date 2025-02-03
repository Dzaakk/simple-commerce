package models

type UpdateReq struct {
	Id          string `json:"id" validate:"required,min=1"`
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type DataRes struct {
	Id          string `json:"id"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Balance     string `json:"balance,omitempty"`
}

type BalanceUpdateReq struct {
	Id         string `json:"id" validate:"required,numeric,min=1"`
	ActionType string `json:"actionType" validate:"required"`
	Balance    string `json:"balance" validate:"required"`
}
type BalanceUpdateRes struct {
	BalanceOld CustomerBalanceRes `json:"oldData"`
	BalanceNew CustomerBalanceRes `json:"newData"`
}

type CustomerBalanceRes struct {
	Id      string `json:"id"`
	Balance string `json:"balance"`
}
