package models

import "Dzaakk/simple-commerce/package/template"

type TCustomers struct {
	Id          int     `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	Password    string  `json:"password"`
	Balance     float32 `json:"balance"`
	template.Base
}
type CustomerReq struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type CustomerBalance struct {
	Id      int     `json:"id"`
	Balance float32 `json:"balance"`
}
type CustomerRes struct {
	Id          string `json:"id"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Balance     string `json:"balance,omitempty"`
}
type BalanceUpdateReq struct {
	Id         string `json:"id"`
	ActionType string `json:"actionType"`
	Balance    string `json:"balance"`
}
type BalanceUpdateRes struct {
	BalanceOld CustomerBalance `json:"oldData"`
	BalanceNew CustomerBalance `json:"newData"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
