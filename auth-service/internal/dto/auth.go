package dto

import uuid "github.com/satori/go.uuid"

type Auth struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email" `
	IsAdmin bool      `json:"isAdmin" `
	Token   string    `json:"token"`
}
