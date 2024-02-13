package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type UserController struct{}

func (controller *UserController) RegisterRoute(router *mux.Router) {
	log.Println("------Initializing User Controller-------")
	userRouter := router.PathPrefix("/api").Subrouter()

	userRouter.HandleFunc("/user", controller.CreateUser).Methods(http.MethodPost)
	userRouter.HandleFunc("/user", controller.GetUsers).Methods(http.MethodGet)
	userRouter.HandleFunc("/user/:id", controller.GetUser).Methods(http.MethodGet)
	userRouter.HandleFunc("/user/:id", controller.DeleteUser).Methods(http.MethodDelete)
	userRouter.HandleFunc("/user/:id", controller.UpdateUser).Methods(http.MethodPut)
}

func (controller *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {}

func (controller *UserController) GetUser(w http.ResponseWriter, r *http.Request) {}

func (controller *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {}

func (controller *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {}

func (controller *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {}
