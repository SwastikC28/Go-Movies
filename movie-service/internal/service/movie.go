package service

import (
	"fmt"
	"log"
	"movie-service/internal/model"
	"shared/datastore"
	"shared/datastore/relationaldb"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type MovieService struct {
	db   *gorm.DB
	repo datastore.Repository
}

func NewMovieService(db *gorm.DB, repo datastore.Repository) *MovieService {
	return &MovieService{
		db:   db,
		repo: repo,
	}
}

func (service *MovieService) Create(newMovie *model.Movie) error {
	//  Creating unit of work.
	uow := relationaldb.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	newMovie.ID = uuid.NewV4()

	// Add newMovie.
	err := service.repo.Add(uow, newMovie)
	if err != nil {
		uow.Rollback()
		return err
	}

	uow.Commit()
	return nil
}

func (service *MovieService) GetAllMovies(Movies *[]model.Movie) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)

	defer uow.Rollback()

	err := service.repo.GetAllRecords(uow, Movies, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}

func (service *MovieService) GetMovie(movie *model.Movie, queryProcessor []datastore.QueryProcessor) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)

	defer uow.Rollback()

	err := service.repo.GetFirst(uow, &movie, queryProcessor)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}

func (service *MovieService) DeleteMovie(id string) error {
	uow := relationaldb.NewUnitOfWork(service.db, false)

	condition := fmt.Sprintf("ID = %s", id)

	defer uow.Rollback()
	err := service.repo.Delete(uow, &model.Movie{}, condition)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}

func (service *MovieService) UpdateMovie(movie *model.Movie) error {
	uow := relationaldb.NewUnitOfWork(service.db, false)
	defer uow.Rollback()

	// Check if the movie exists
	err := service.repo.GetFirst(uow, &model.Movie{}, []datastore.QueryProcessor{datastore.Filter("id = ?", movie.ID)})
	if err != nil {
		log.Println(err)
		return err
	}

	// Update Movie
	err = service.repo.Update(uow, &movie)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}
