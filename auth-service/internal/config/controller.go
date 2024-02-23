package config

import (
	"auth-service/internal/controller"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"shared/datastore"

	"github.com/gorilla/mux"
)

type Controller interface {
	RegisterRoute(router *mux.Router)
}

func RegisterAuthRoutes(app *App) {
	defer app.WG.Done()

	authService := service.NewAuthService(app.DB, &repository.UserRepository{
		GormRepository: *datastore.NewGormRepository(),
	})

	authController := controller.NewAuthController(authService)

	authController.RegisterRoutes(app.Router)
}

func RegisterRoutes(app *App) {
	app.WG.Add(1)

	go RegisterAuthRoutes(app)

	app.WG.Wait()
}
