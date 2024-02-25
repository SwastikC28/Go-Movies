package service

import (
	"fmt"
	"log"
	"rental-service/internal/model"
	"shared/datastore"
	"shared/datastore/relationaldb"

	"github.com/jinzhu/gorm"
)

type RentalService struct {
	db   *gorm.DB
	repo datastore.Repository
}

func NewRentalService(db *gorm.DB, repo datastore.Repository) *RentalService {
	return &RentalService{
		db:   db,
		repo: repo,
	}
}

func (service *RentalService) Create(newRental *model.Rental) error {
	//  Creating unit of work.
	uow := relationaldb.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	// Add newRental.
	err := service.repo.Add(uow, newRental)
	if err != nil {
		uow.Rollback()
		return err
	}

	uow.Commit()
	return nil
}

func (service *RentalService) GetAllRentals(Rentals *[]model.Rental) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)

	defer uow.Rollback()

	err := service.repo.GetAllRecords(uow, Rentals, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}

func (service *RentalService) GetRental(Rental *model.Rental, queryProcessor []datastore.QueryProcessor) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)

	defer uow.Rollback()

	err := service.repo.GetFirst(uow, &Rental, queryProcessor)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}

func (service *RentalService) DeleteRental(id string) error {
	uow := relationaldb.NewUnitOfWork(service.db, false)

	condition := fmt.Sprintf("ID = %s", id)

	defer uow.Rollback()
	err := service.repo.Delete(uow, &model.Rental{}, condition)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}
