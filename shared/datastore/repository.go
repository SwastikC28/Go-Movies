package datastore

import (
	"shared/datastore/relationaldb"

	"github.com/jinzhu/gorm"
)

type Repository interface {
	GetAllRecords()
	GetFirst()
	Add()
	Update()
	Save()
}

type GormRepository struct{}

func NewGormRepository() *GormRepository {
	return &GormRepository{}
}

// Repo Methods
// filter
// select

func (repo *GormRepository) Filter(condition string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Where(condition, args...)
		return db, nil
	}
}

func (repo *GormRepository) Select(condition string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Select(condition, args...)
		return db, nil
	}
}

func (repo *GormRepository) GetAllRecords(uow *relationaldb.UnitOfWork, out interface{}, queryProcessors []QueryProcessor) error {
	db := uow.DB
	var err error
	for _, query := range queryProcessors {
		db, err = query(db, out)
		if err != nil {
			return err
		}
	}

	return db.Debug().Find(out).Error
}

func (repo *GormRepository) GetFirst(uow *relationaldb.UnitOfWork, out interface{}, queryProcessors []QueryProcessor) error {
	db := uow.DB
	var err error
	for _, query := range queryProcessors {
		db, err = query(db, out)
		if err != nil {
			return err
		}
	}

	return db.Debug().First(out).Error

}

func (repo *GormRepository) Add(uow *relationaldb.UnitOfWork, out interface{}) error {
	db := uow.DB
	return db.Create(out).Debug().Error
}

func (repo *GormRepository) Update(uow *relationaldb.UnitOfWork, out interface{}) error {
	db := uow.DB
	return db.Model(out).Update(out).Error
}

func (repo *GormRepository) Save(uow *relationaldb.UnitOfWork, out interface{}) error {
	db := uow.DB
	return db.Debug().Save(out).Error
}
