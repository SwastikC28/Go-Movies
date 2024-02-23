package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(100)"`
	IsAdmin  bool   `json:"isAdmin" gorm:"type:boolean;default:false"`
	Email    string `json:"email" gorm:"type:varchar(100);unique"`
	Password string `json:"password" gorm:"type:varchar(100);"`
}
