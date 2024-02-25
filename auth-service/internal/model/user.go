package model

import (
	"shared/pkg/model"
)

type User struct {
	model.Base
	Name     string `json:"name" gorm:"type:varchar(100)"`
	IsAdmin  bool   `json:"isAdmin" gorm:"type:boolean;default:false"`
	Email    string `json:"email" gorm:"type:varchar(100);unique"`
	Password string `json:"password" gorm:"type:string;"`
}
