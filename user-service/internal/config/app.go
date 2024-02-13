package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	Name   string
	WG     *sync.WaitGroup
	Server http.Server
	Router *mux.Router
	DB     *gorm.DB
}

func NewApp(name string, db *gorm.DB, wg *sync.WaitGroup) *App {
	return &App{
		Name: name,
		WG:   wg,
		DB:   db,
	}
}

func (app *App) InitializeRouter() {
	app.Router = mux.NewRouter().StrictSlash(true)
	app.Router = app.Router.PathPrefix("/").Subrouter()

	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Server is Running"))
	})
}

func (app *App) InitializeServer() {
	app.Server = http.Server{
		Addr:    "0.0.0.0:4000",
		Handler: app.Router,
	}
}

func (app *App) Init() {
	app.InitializeRouter()
	app.InitializeServer()
}

func (app *App) StartServer() error {
	err := app.Server.ListenAndServe()
	if err != nil {
		return err
	}

	log.Println("Server Started on PORT 4000")
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

	// close db connection
	app.DB.Close()

	log.Println("Server shutdown successfully.")
}
