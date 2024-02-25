package model

import "shared/pkg/model"

type Auth struct {
	model.Base
	Name    string `json:"name" gorm:"type:varchar(100)"`
	Email   string `json:"email" gorm:"type:varchar(100)"`
	IsAdmin bool   `json:"isAdmin" gorm:"type:boolean;default:false"`
}
