package controller

import (
	"fmt"
	"movie-service/internal/model"
	"movie-service/internal/service"
	"net/http"
	"shared/datastore"
	"shared/middleware"
	"shared/pkg/web"
	"shared/security"

	"github.com/gorilla/mux"
)

type MovieController struct {
	service *service.MovieService
}

func (controller *MovieController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----Movie Controller Registered-----")
	movieRouter := router.PathPrefix("/movie").Subrouter()

	movieRouter.Use(middleware.ReqLogger)

	movieRouter.HandleFunc("", web.AccessGuard(controller.createMovie)).Methods(http.MethodPost)
	movieRouter.HandleFunc("", controller.getMovies).Methods(http.MethodGet)
	movieRouter.HandleFunc("/{id}", controller.getMovieById).Methods(http.MethodGet)
	movieRouter.HandleFunc("/{id}", web.AccessGuard(controller.deleteMovieById)).Methods(http.MethodDelete)
}

func NewMovieController(service *service.MovieService) *MovieController {
	return &MovieController{
		service: service,
	}
}

func (controller *MovieController) createMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	web.UnmarshalJSON(r, &movie)

	token := security.TokenFromContext(r.Context())

	if !token.IsAdmin {
		web.RespondJSON(w, http.StatusUnauthorized, "User unauthorized to access this route")
		return
	}

	// Assign Created By as the user's id from token
	movie.CreatedBy = token.ID

	err := controller.service.Create(&movie)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, movie)
}

func (controller *MovieController) getMovies(w http.ResponseWriter, r *http.Request) {
	var movies []model.Movie

	err := controller.service.GetAllMovies(&movies)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, movies)
}

func (controller *MovieController) getMovieById(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie

	vars := mux.Vars(r)

	id := vars["id"]

	queryProcessor := []datastore.QueryProcessor{}
	queryProcessor = append(queryProcessor, datastore.Filter("ID = ?", (id)))

	err := controller.service.GetMovie(&movie, queryProcessor)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, movie)
}

func (controller *MovieController) deleteMovieById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	err := controller.service.DeleteMovie(id)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, "Movie Deleted Successfully.")
}
