package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:varchar(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
