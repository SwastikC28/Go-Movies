package controller

import (
	"auth-service/internal/dto"
	"auth-service/internal/model"
	"auth-service/internal/service"
	"fmt"
	"net/http"
	"shared/auth"
	"shared/middleware"
	"shared/utils/web"

	"github.com/gorilla/mux"
)

type AuthController struct {
	service *service.AuthService
}

func (controller *AuthController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----Auth Controller Registered-----")
	authRouter := router.PathPrefix("/auth").Subrouter()

	authRouter.Use(middleware.ReqLogger)

	authRouter.HandleFunc("/login", controller.login).Methods(http.MethodPost)
	authRouter.HandleFunc("/register", controller.register).Methods(http.MethodPost)

}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

func (controller *AuthController) login(w http.ResponseWriter, r *http.Request) {
	// User Info
	var user model.User
	web.UnmarshalJSON(r, &user)

	// Match Passwords
	err := controller.service.MatchPassword(&user)
	if err != nil {
		web.RespondJSON(w, http.StatusUnauthorized, "Invalid Email or Password")
		return
	}

	// Get Token
	userClaim := auth.Claims{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}

	token, err := auth.SignJWT(userClaim)

	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, "error while creating the token")
		return
	}

	authDTO := dto.Auth{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		Token:   token,
	}

	web.RespondJSON(w, http.StatusOK, authDTO)
}

func (controller *AuthController) register(w http.ResponseWriter, r *http.Request) {
	// User Info
	var user model.User
	web.UnmarshalJSON(r, &user)

	err := controller.service.Create(&user)
	if err != nil {
		web.RespondJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Get Token
	userClaim := auth.Claims{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}

	token, err := auth.SignJWT(userClaim)

	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, "error while creating the token")
		return
	}

	authDTO := dto.Auth{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		Token:   token,
	}

	web.RespondJSON(w, http.StatusOK, authDTO)
}
