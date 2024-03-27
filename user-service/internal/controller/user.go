package controller

import (
	"fmt"
	"net/http"
	"shared/datastore"
	"shared/middleware"
	"shared/pkg/web"
	"shared/security"
	"user-service/internal/model"
	"user-service/internal/service"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type UserController struct {
	service *service.UserService
}

func (controller *UserController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----User Controller Registered-----")
	userRouter := router.PathPrefix("/user").Subrouter()

	userRouter.Use(middleware.ReqLogger)

	userRouter.HandleFunc("", web.AccessGuard(controller.createUser, true)).Methods(http.MethodPost)
	userRouter.HandleFunc("", web.AccessGuard(controller.getUsers, true)).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", web.AccessGuard(controller.getUserById, false)).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", web.AccessGuard(controller.deleteUserById, false)).Methods(http.MethodDelete)
	userRouter.HandleFunc("/{id}", web.AccessGuard(controller.updateUser, false)).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id}/change-password", web.AccessGuard(controller.updatePassword, false)).Methods(http.MethodPut)
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (controller *UserController) createUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	web.UnmarshalJSON(r, &user)

	hashedPassword, err := web.EncryptPassword(user.Password)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = string(hashedPassword)

	// Create User
	err = controller.service.Create(&user)
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

func (controller *UserController) updateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	web.UnmarshalJSON(r, &user)
	vars := mux.Vars(r)

	id := vars["id"]

	user.ID = uuid.FromStringOrNil(id)

	token := security.TokenFromContext(r.Context())

	if token.ID.String() != id && !token.IsAdmin {
		web.RespondJSON(w, http.StatusUnauthorized, "User unauthorized to access this route")
		return
	}

	// Should not be able to change password
	user.Password = ""

	err := controller.service.UpdateUser(&user)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, "User Updated Successfully.")
}

func (controller *UserController) updatePassword(w http.ResponseWriter, r *http.Request) {
	var user = struct {
		Email           string `json:"email"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}{}

	web.UnmarshalJSON(r, &user)
	vars := mux.Vars(r)

	id := vars["id"]

	token := security.TokenFromContext(r.Context())

	if token.ID.String() != id && !token.IsAdmin {
		web.RespondJSON(w, http.StatusUnauthorized, "User unauthorized to access this route")
		return
	}

	var existingUser model.User
	err := controller.service.GetUser(&existingUser, []datastore.QueryProcessor{datastore.Filter("email = ?", user.Email)})
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = web.ComparePassword(user.CurrentPassword, []byte(existingUser.Password))
	if err != nil {
		web.RespondJSON(w, http.StatusUnauthorized, "Incorrect current password. Please try again.")
		return
	}

	err = controller.service.ChangedPassword(&existingUser, user.NewPassword)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, "User Password Changed Successfully.")
}
