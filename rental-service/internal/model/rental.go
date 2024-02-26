package model

import (
	"shared/pkg/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type AvailabilityStatus string

const (
	Paid   AvailabilityStatus = "paid"
	Unpaid AvailabilityStatus = "unpaid"
)

type Rental struct {
	model.Base
	RentalDate time.Time
	DueDate    time.Time          `json:"dueDate"`
	ReturnDate *time.Time         `json:"returnDate"`
	Status     AvailabilityStatus `json:"status" gorm:"default:'unpaid'"`
	LateFee    uint               `json:"lateFee"`
	MovieId    uuid.UUID          `json:"movieId" gorm:"type:varchar(36);foreignKey"`
	UserId     uuid.UUID          `json:"userId" gorm:"type:varchar(36);foreignKey"`
}
