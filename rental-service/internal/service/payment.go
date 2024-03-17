package service

import (
	"errors"
	"rental-service/internal/model"
	"shared/datastore"
	"shared/datastore/relationaldb"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type PaymentService struct {
	db   *gorm.DB
	repo datastore.Repository
}

func NewPaymentService(db *gorm.DB, repo datastore.Repository) *PaymentService {
	return &PaymentService{
		db:   db,
		repo: repo,
	}
}

func (service *PaymentService) GetRentalFees(rental *model.Rental) error {
	//  Creating unit of work.
	uow := relationaldb.NewUnitOfWork(service.db, false)
	defer uow.Rollback()

	err := service.repo.GetFirst(uow, rental, []datastore.QueryProcessor{datastore.Filter("id = ?", rental.ID)})
	if err != nil {
		return err
	}

	if rental.RentalFee == 0 {
		return errors.New("please return movie before paying")
	}

	if rental.Status == "paid" {
		return errors.New("rental fees already paid")
	}

	return nil
}

func (service *PaymentService) SavePayment(payment model.Payment) error {
	uow := relationaldb.NewUnitOfWork(service.db, false)
	defer uow.Rollback()

	// Create Entry
	payment.ID = uuid.NewV4()
	payment.CreatedAt = time.Now()

	// Add in Payment's Database
	err := service.repo.Add(uow, &payment)
	if err != nil {
		return err
	}

	// Update Rental Row to true
	var rental model.Rental
	err = service.repo.GetFirst(uow, &rental, []datastore.QueryProcessor{datastore.Filter("id = ?", payment.RentalId)})
	if err != nil {
		return err
	}

	if rental.Status == "paid" {
		return errors.New("rental fees already paid")
	}

	// Update Status to Paid
	rental.Status = "paid"
	err = service.repo.Save(uow, &rental)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}
