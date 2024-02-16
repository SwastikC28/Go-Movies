package controller

import (
	"fmt"
	"net/http"
	"shared/utils/web"

	"github.com/gorilla/mux"
)

type UserController struct {
}

func (controller *UserController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----User Controller Registered-----")
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", controller.createUser).Methods(http.MethodPost)
	userRouter.HandleFunc("/", controller.getUser).Methods(http.MethodGet)
}

func (controller *UserController) createUser(w http.ResponseWriter, r *http.Request) {
	web.RespondJSON(w, http.StatusOK, "Created User Successfully")
}

func (controller *UserController) getUser(w http.ResponseWriter, r *http.Request) {
	web.RespondJSON(w, http.StatusOK, "Get User Successfully")
}
