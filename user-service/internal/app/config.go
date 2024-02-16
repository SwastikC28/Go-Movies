package app

import (
	"fmt"
	"user-service/internal/config"
	"user-service/internal/controller"
	"user-service/internal/model"

	"github.com/gorilla/mux"
)

type Controller interface {
	RegisterRoutes(router *mux.Router)
}

func RegisterRoute(app *config.App) {
	userController := &controller.UserController{}

	userController.RegisterRoutes(app.Router)
}

func TableMigration(app *config.App) {
	fmt.Println("-----User Table Migration-----")
	app.DB.AutoMigrate(&model.User{})
}
