package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string
	Type     string
	Email    string
	Password string
}
