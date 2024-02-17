package datastore

import "github.com/jinzhu/gorm"

type QueryProcessor func(db *gorm.DB, out interface{}) (*gorm.DB, error)

func Filter(condition string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Where(condition, args...)
		return db, nil
	}
}

func OrderBy(condition string) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Order(condition)
		return db, nil
	}
}
