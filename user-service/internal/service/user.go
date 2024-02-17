package service

import (
	"fmt"
	"log"
	"shared/datastore"
	"shared/datastore/relationaldb"
	"user-service/internal/model"

	"github.com/jinzhu/gorm"
)

type UserService struct {
	db   *gorm.DB
	repo datastore.Repository
}

func NewUserService(db *gorm.DB, repo datastore.Repository) *UserService {
	return &UserService{
		db:   db,
		repo: repo,
	}
}

func (service *UserService) Create(newUser *model.User) error {
	//  Creating unit of work.
	uow := relationaldb.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	// Add newUser.
	err := service.repo.Add(uow, newUser)
	if err != nil {
		uow.Rollback()
		return err
	}

	uow.Commit()
	return nil
}

func (service *UserService) GetAllUsers(users *[]model.User) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)

	defer uow.Rollback()

	err := service.repo.GetAllRecords(uow, users, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}

func (service *UserService) GetUser(user *model.User, queryProcessor []datastore.QueryProcessor) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)

	defer uow.Rollback()

	err := service.repo.GetFirst(uow, &user, queryProcessor)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}

func (service *UserService) DeleteUser(id string) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)

	condition := fmt.Sprintf("ID = %s", id)

	defer uow.Rollback()
	err := service.repo.Delete(uow, &model.User{}, condition)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}
