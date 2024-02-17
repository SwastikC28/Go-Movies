package relationaldb

import "github.com/jinzhu/gorm"

type UnitOfWork struct {
	DB        *gorm.DB
	Committed bool
	Readonly  bool
}

func NewUnitOfWork(db *gorm.DB, readonly bool) *UnitOfWork {
	commit := false
	if readonly {
		return &UnitOfWork{
			DB:        db.New(),
			Committed: commit,
			Readonly:  readonly,
		}
	}

	return &UnitOfWork{
		DB:        db.New().Begin(),
		Committed: commit,
		Readonly:  readonly,
	}
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
