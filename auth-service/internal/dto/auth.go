package dto

type Auth struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email" `
	IsAdmin bool   `json:"isAdmin" `
	Token   string `json:"token"`
}
