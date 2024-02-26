package repository

import (
	"shared/datastore"
)

type RentalRepository struct {
	datastore.GormRepository
}
