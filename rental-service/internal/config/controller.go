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

func RegisterMovieRoutes(app *App) {
	defer app.WG.Done()

	rentalService := service.NewRentalService(app.DB, &repository.MovieRepository{
		GormRepository: *datastore.NewGormRepository(),
	})

	movieController := controller.NewRentalController(rentalService)

	movieController.RegisterRoutes(app.Router)
}

func TableMigration(app *App) {
	fmt.Println("-----Movie Table Migration-----")
	err := app.DB.AutoMigrate(&model.Movie{}).Error
	if err != nil {
		log.Println(err)
	}
}

func RegisterRoutes(app *App) {
	app.WG.Add(1)

	go RegisterMovieRoutes(app)

	app.WG.Wait()
}
