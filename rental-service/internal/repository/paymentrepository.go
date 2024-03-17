package repository

import "shared/datastore"

type PaymentRepository struct {
	datastore.GormRepository
}
