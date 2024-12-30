package models

type ReqCreate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ReqUpdate struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ResCreate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ResData struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Balance  string `json:"balance"`
}
