package app

import (
	"log"
	"net/http"
	"sync"

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
	app.Router = app.Router.PathPrefix("/api/user").Subrouter()

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
