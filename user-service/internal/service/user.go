package service

import (
	"shared/datastore"

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

func (service *UserService) Create(out interface{}) error {
	return nil
}
