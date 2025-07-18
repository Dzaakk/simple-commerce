package model

import (
	"Dzaakk/simple-commerce/package/template"
	"fmt"
)

type TSeller struct {
	ID                int64   `json:"id"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	PhoneNumber       string  `json:"phone_number"`
	StoreName         string  `json:"store_name"`
	Address           string  `json:"address"`
	Balance           float64 `json:"balance"`
	Status            int     `json:"status"`
	ProfilePicture    string  `json:"profile_picture"`
	BankAccountName   string  `json:"bank_account_name"`
	BankAccountNumber string  `json:"bank_account_number"`
	BankName          string  `json:"bank_name"`
	template.Base
}

func (s *TSeller) ToResponse() SellerRes {
	return SellerRes{
		ID:             fmt.Sprintf("%d", s.ID),
		Username:       s.Username,
		Email:          s.Email,
		Balance:        fmt.Sprintf("%.2f", s.Balance),
		StoreName:      s.StoreName,
		ProfilePicture: s.ProfilePicture,
		Address:        s.Address,
	}
}
