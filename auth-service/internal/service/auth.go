package service

import (
	"auth-service/internal/model"
	"shared/datastore"
	"shared/datastore/relationaldb"
	"shared/utils/web"

	"github.com/jinzhu/gorm"
)

type AuthService struct {
	db   *gorm.DB
	repo datastore.Repository
}

func NewAuthService(db *gorm.DB, repo datastore.Repository) *AuthService {
	return &AuthService{
		db:   db,
		repo: repo,
	}
}

func (service *AuthService) Create(newUser *model.User) error {
	//  Creating unit of work.
	uow := relationaldb.NewUnitOfWork(service.db, false)

	// Get User if exist throw error

	defer uow.Rollback()

	// encrypt password
	hashedPassword, err := web.EncryptPassword(newUser.Password)
	if err != nil {
		return err
	}

	newUser.Password = string(hashedPassword)

	// Add new user.
	err = service.repo.Add(uow, newUser)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

func (service *AuthService) MatchPassword(userData *model.User) error {
	uow := relationaldb.NewUnitOfWork(service.db, true)

	// Get User Data
	var user = &model.User{}

	err := service.repo.GetFirst(uow, user, []datastore.QueryProcessor{datastore.Filter("email = ?", userData.Email)})
	if err != nil {
		return err
	}

	err = web.ComparePassword(userData.Password, []byte(user.Password))
	if err != nil {
		return err
	}

	userData.ID = user.ID
	userData.Name = user.Name

	return nil
}
