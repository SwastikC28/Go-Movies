package config

import (
	"auth-service/internal/controller"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"fmt"
	"os"
	"shared/datastore"
	"shared/pkg/event/monitor"
	"shared/pkg/event/publisher"

	"github.com/gorilla/mux"
)

var eventHandler monitor.EventHandler

type Controller interface {
	RegisterRoute(router *mux.Router)
}

func RegisterAuthRoutes(app *App) {
	defer app.WG.Done()

	authService := service.NewAuthService(app.DB, &repository.UserRepository{
		GormRepository: *datastore.NewGormRepository(),
	})

	exchangeKey := os.Getenv("RABBITMQ_EXCHANGE_KEY")

	dispatcher, err := publisher.NewDispatchHandler(exchangeKey)
	if err != nil {
		fmt.Println("There was some problem creating dispatcher", err)
	}

	// Assign Event with Implementation
	eventHandler = dispatcher

	authController := controller.NewAuthController(authService, dispatcher)
	authController.RegisterRoutes(app.Router)
}

func RegisterRoutes(app *App) {
	app.WG.Add(1)

	go RegisterAuthRoutes(app)

	app.WG.Wait()
}
