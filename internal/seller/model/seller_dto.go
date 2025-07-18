package model

type ReqCreate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ReqUpdate struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ResCreate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type SellerRes struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Balance        string `json:"balance"`
	StoreName      string `json:"store_name"`
	ProfilePicture string `json:"profile_picture"`
	Address        string `json:"address"`
}
