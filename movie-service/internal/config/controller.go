package config

import (
	"fmt"
	"log"
	"movie-service/internal/controller"
	"movie-service/internal/model"
	"movie-service/internal/repository"
	"movie-service/internal/service"
	"shared/datastore"

	"github.com/gorilla/mux"
)

type Controller interface {
	RegisterRoute(router *mux.Router)
}

func RegisterMovieRoutes(app *App) {
	defer app.WG.Done()

	movieService := service.NewMovieService(app.DB, &repository.MovieRepository{
		GormRepository: *datastore.NewGormRepository(),
	})

	movieController := controller.NewMovieController(movieService)

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
