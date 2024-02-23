package repository

import (
	"shared/datastore"
)

type UserRepository struct {
	datastore.GormRepository
}
