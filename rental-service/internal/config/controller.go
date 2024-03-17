package config

import (
	"fmt"
	"log"
	"rental-service/internal/controller"
	"rental-service/internal/model"

	"rental-service/internal/repository"
	"rental-service/internal/service"
	"shared/datastore"

	"github.com/gorilla/mux"
)

type Controller interface {
	RegisterRoute(router *mux.Router)
}

func RegisterRentalRoutes(app *App) {
	defer app.WG.Done()

	rentalService := service.NewRentalService(app.DB, &repository.RentalRepository{
		GormRepository: *datastore.NewGormRepository(),
	})

	rentalController := controller.NewRentalController(rentalService)

	rentalController.RegisterRoutes(app.Router)
}

func RegisterPaymentRoutes(app *App) {
	defer app.WG.Done()

	paymentService := service.NewPaymentService(app.DB, &repository.RentalRepository{
		GormRepository: *datastore.NewGormRepository(),
	})

	paymentController := controller.NewPaymentController(paymentService)

	paymentController.RegisterRoutes(app.Router)
}

func TableMigration(app *App) {
	fmt.Println("-----Rental MS Tables Migration-----")

	err := app.DB.AutoMigrate(&model.Rental{}).Error
	if err != nil {
		log.Println(err)
	}

	err = app.DB.AutoMigrate(&model.Order{}).Error
	if err != nil {
		log.Println(err)
	}

	err = app.DB.AutoMigrate(&model.Payment{}).Error
	if err != nil {
		log.Println(err)
	}

	if err := app.DB.Model(&model.Rental{}).
		AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("movie_id", "movies(id)", "RESTRICT", "RESTRICT"); err != nil {
		fmt.Println(err)
	}
}

func RegisterRoutes(app *App) {
	app.WG.Add(2)

	go RegisterRentalRoutes(app)
	go RegisterPaymentRoutes(app)

	app.WG.Wait()
}
