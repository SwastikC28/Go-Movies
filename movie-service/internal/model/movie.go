package model

import (
	"shared/pkg/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Movie struct {
	model.Base
	Title           string     `json:"title" gorm:"type:varchar(100)"`
	Genre           string     `json:"genre" gorm:"type:varchar(30)"`
	Release_Date    *time.Time `json:"releaseDate"`
	Director        string     `json:"director" gorm:"type:varchar(50);"`
	Description     string     `json:"description" gorm:"type:varchar(100);"`
	Inventory_Count uint       `json:"inventoryCount" gorm:"type:integer;default:0;"`
	CreatedBy       uuid.UUID  `json:"createdBy" gorm:"type:varchar(36);foreignKey"`
}
