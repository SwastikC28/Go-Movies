package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Movie struct {
	gorm.Model
	Title           string     `json:"name" gorm:"type:varchar(100)"`
	Genre           bool       `json:"genre" gorm:"type:varchar(10)"`
	Release_Date    *time.Time `json:"releaseDate"`
	Director        string     `json:"director" gorm:"type:varchar(20);"`
	Description     string     `json:"description" gorm:"type:varchar(100);"`
	Inventory_Count uint       `json:"inventoryCount" gorm:"type:integer,default:0;"`
}
