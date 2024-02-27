package service

import (
	"fmt"
	"log"
	"rental-service/internal/model"
	"shared/datastore"
	"shared/datastore/relationaldb"
	"shared/pkg/web"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type RentalService struct {
	db   *gorm.DB
	repo datastore.Repository
	gorm.Association
	associations []string
}

func NewRentalService(db *gorm.DB, repo datastore.Repository) *RentalService {
	return &RentalService{
		db:           db,
		repo:         repo,
		associations: []string{"User", "Movie"},
	}
}

func (service *RentalService) Create(newRental *model.Rental) error {
	//  Creating unit of work.
	uow := relationaldb.NewUnitOfWork(service.db, false)
	defer uow.Rollback()

	// Check if movie inventory is greater than 0
	var movie model.Movie
	movie.ID = newRental.MovieId
	err := service.GetMovie(&movie)
	if err != nil {
		return gorm.ErrRecordNotFound
	}

	if movie.Inventory_Count <= 0 {
		return gorm.ErrRecordNotFound
	}

	// Add newRental.
	newRental.ID = uuid.NewV4()

	err = service.repo.Add(uow, newRental)
	if err != nil {
		return err
	}

	// Decrement movie count & save
	movie.Inventory_Count = movie.Inventory_Count - 1
	service.repo.Save(uow, movie)

	uow.Commit()
	return nil
}

func (service *RentalService) GetAllRentals(rentals *[]model.Rental, includes []string, queryProcessors []datastore.QueryProcessor) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)
	defer uow.Rollback()

	// Get matched associations
	includes = web.GetMatchedAssociations(includes, service.associations)

	// Append Preload QueryProcessor
	queryProcessors = append(queryProcessors, datastore.Preload(includes))

	err := service.repo.GetAllRecords(uow, rentals, queryProcessors)
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

func (service *RentalService) GetMovie(movie *model.Movie) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)
	defer uow.Rollback()

	queryProcessor := []datastore.QueryProcessor{}
	queryProcessor = append(queryProcessor, datastore.Filter("id= ?", movie.ID))

	err := service.repo.GetFirst(uow, movie, queryProcessor)

	if err != nil {
		return err
	}

	return nil
}
