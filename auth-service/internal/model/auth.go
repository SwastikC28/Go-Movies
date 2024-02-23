package model

import "github.com/jinzhu/gorm"

type Auth struct {
	gorm.Model
	Name    string `json:"name" gorm:"type:varchar(100)"`
	Email   string `json:"email" gorm:"type:varchar(100)"`
	IsAdmin bool   `json:"isAdmin" gorm:"type:boolean;default:false"`
}
