package relationaldb

import "github.com/jinzhu/gorm"

type UnitOfWork struct {
	DB        *gorm.DB
	Committed bool
	Readonly  bool
}

func NewUnitOfWork(db *gorm.DB, readonly bool) *UnitOfWork {
	uow := &UnitOfWork{
		DB:        db,
		Committed: false,
		Readonly:  true,
	}

	if !readonly {
		uow.Readonly = false
	}

	return uow
}

func (uow *UnitOfWork) Commit() {
	if !uow.Readonly && !uow.Committed {
		uow.Committed = true
		uow.DB.Commit()
	}
}

func (uow *UnitOfWork) Rollback() {
	if !uow.Committed && !uow.Readonly {
		uow.DB.Rollback()
	}
}
