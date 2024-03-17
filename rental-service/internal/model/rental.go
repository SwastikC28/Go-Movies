package model

import (
	"shared/pkg/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Rental struct {
	model.Base
	RentalDate time.Time
	DueDate    time.Time  `json:"dueDate"`
	ReturnDate *time.Time `json:"returnDate"`
	Status     string     `json:"status" gorm:"default:'unpaid'"`
	LateFee    float64    `json:"lateFee"`
	RentalFee  int        `json:"rentalFee"`
	MovieId    uuid.UUID  `json:"movieId" gorm:"type:varchar(36)"`
	Movie      Movie      `gorm:"foreignKey:MovieId"`
	UserId     uuid.UUID  `json:"userId" gorm:"type:varchar(36);foreignKey"`
	User       User       `gorm:"foreignKey:UserId"`
}
