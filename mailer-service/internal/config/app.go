package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type App struct {
	Name   string
	AMQP   *amqp.Connection
	Server http.Server
}

func NewApp(name string, AMQP *amqp.Connection) *App {
	return &App{
		Name: name,
		AMQP: AMQP,
	}
}

func (app *App) InitializeServer() {
	app.Server = http.Server{
		Addr: "0.0.0.0:80",
	}
}

func (app *App) Init() {
	app.InitializeServer()
}

func (app *App) StartServer() error {
	err := app.Server.ListenAndServe()
	if err != nil {
		return err
	}

	log.Println("Server Started on PORT 80")
	return nil
}

func (app *App) StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	err := app.Server.Shutdown(ctx)
	if err != nil {
		fmt.Println("Fail to stop server.")
		return
	}

	log.Println("Server shutdown successfully.")
}
