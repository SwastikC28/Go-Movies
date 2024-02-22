package controller

import (
	"fmt"
	"movie-service/internal/model"
	"movie-service/internal/service"
	"net/http"
	"shared/datastore"
	"shared/middleware"
	"shared/utils/web"

	"github.com/gorilla/mux"
)

type MovieController struct {
	service *service.MovieService
}

func (controller *MovieController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----Movie Controller Registered-----")
	movieRouter := router.PathPrefix("/movie").Subrouter()

	movieRouter.Use(middleware.ReqLogger)

	movieRouter.HandleFunc("", controller.createMovie).Methods(http.MethodPost)
	movieRouter.HandleFunc("", controller.getMovies).Methods(http.MethodGet)
	movieRouter.HandleFunc("/{id}", controller.getMovieById).Methods(http.MethodGet)
	movieRouter.HandleFunc("/{id}", controller.deleteMovieById).Methods(http.MethodDelete)
}

func NewMovieController(service *service.MovieService) *MovieController {
	return &MovieController{
		service: service,
	}
}

func (controller *MovieController) createMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	web.UnmarshalJSON(r, &movie)

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
	fmt.Println("ID", id)

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
