package models

type CustomerReq struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
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
