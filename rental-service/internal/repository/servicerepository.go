package repository

import (
	"shared/datastore"
)

type MovieRepository struct {
	datastore.GormRepository
}
