package service

import (
	"log"
	"shared/datastore"
	"shared/datastore/relationaldb"
	"user-service/internal/model"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
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

func (service *UserService) GetDB() *gorm.DB {
	return service.db.New()
}

func (service *UserService) Create(newUser *model.User) error {
	//  Creating unit of work.
	uow := relationaldb.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	// New ID
	newUser.ID = uuid.NewV4()
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
	uow := relationaldb.NewUnitOfWork(service.db, false)

	defer uow.Rollback()

	// Check if the user exists
	user := model.User{}

	err := service.repo.GetFirst(uow, &user, []datastore.QueryProcessor{datastore.Filter("id =?", uuid.FromStringOrNil(id))})
	if err != nil {
		log.Println(err)
		return err
	}

	// Delete User
	err = service.repo.Delete(uow, &user, "")
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}

func (service *UserService) UpdateUser(user *model.User) error {
	uow := relationaldb.NewUnitOfWork(service.db, false)
	defer uow.Rollback()

	// Check if the user exists
	err := service.repo.GetFirst(uow, &model.User{}, []datastore.QueryProcessor{datastore.Filter("id =?", user.ID)})
	if err != nil {
		log.Println(err)
		return err
	}

	// Update User
	err = service.repo.Update(uow, &user)
	if err != nil {
		log.Println(err)
		return err
	}

	uow.Commit()
	return nil
}
