package controller

import (
	"fmt"
	"net/http"
	"shared/datastore"
	"shared/middleware"
	"shared/pkg/web"
	"user-service/internal/model"
	"user-service/internal/service"

	"github.com/gorilla/mux"
)

type UserController struct {
	service *service.UserService
}

func (controller *UserController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----User Controller Registered-----")
	userRouter := router.PathPrefix("/user").Subrouter()

	userRouter.Use(middleware.ReqLogger)

	userRouter.HandleFunc("", controller.createUser).Methods(http.MethodPost)
	userRouter.HandleFunc("", controller.getUsers).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", controller.getUserById).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", controller.deleteUserById).Methods(http.MethodDelete)
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (controller *UserController) createUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	web.UnmarshalJSON(r, &user)

	err := controller.service.Create(&user)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, user)
}

func (controller *UserController) getUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User

	err := controller.service.GetAllUsers(&users)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, users)
}

func (controller *UserController) getUserById(w http.ResponseWriter, r *http.Request) {
	var user model.User

	vars := mux.Vars(r)

	id := vars["id"]

	queryProcessor := []datastore.QueryProcessor{}
	queryProcessor = append(queryProcessor, datastore.Filter("ID = ?", (id)))

	err := controller.service.GetUser(&user, queryProcessor)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, user)
}

func (controller *UserController) deleteUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	err := controller.service.DeleteUser(id)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, "User Deleted Successfully.")
}
