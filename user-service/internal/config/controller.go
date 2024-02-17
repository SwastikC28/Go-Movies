package config

import (
	"fmt"
	"log"
	"shared/datastore"
	"user-service/internal/controller"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/service"

	"github.com/gorilla/mux"
)

type Controller interface {
	RegisterRoute(router *mux.Router)
}

func RegisterUserRoutes(app *App) {
	defer app.WG.Done()

	userService := service.NewUserService(app.DB, &repository.UserRepository{
		GormRepository: *datastore.NewGormRepository(),
	})

	userController := controller.NewUserController(userService)

	userController.RegisterRoutes(app.Router)
}

func TableMigration(app *App) {
	fmt.Println("-----User Table Migration-----")
	err := app.DB.AutoMigrate(&model.User{}).Error
	if err != nil {
		log.Println(err)
	}
}

func RegisterRoutes(app *App) {
	app.WG.Add(1)

	go RegisterUserRoutes(app)

	app.WG.Wait()
}
