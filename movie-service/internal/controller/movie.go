package controller

import (
	"fmt"
	"movie-service/internal/model"
	"movie-service/internal/service"
	"net/http"
	"shared/datastore"
	"shared/middleware"
	"shared/pkg/cloudinary"
	"shared/pkg/web"
	"shared/security"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type MovieController struct {
	service *service.MovieService
}

func (controller *MovieController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----Movie Controller Registered-----")
	movieRouter := router.PathPrefix("/movie").Subrouter()

	movieRouter.Use(middleware.ReqLogger)

	movieRouter.HandleFunc("", web.AccessGuard(controller.createMovie, true)).Methods(http.MethodPost)
	movieRouter.HandleFunc("", controller.getMovies).Methods(http.MethodGet)
	movieRouter.HandleFunc("/{id}", controller.getMovieById).Methods(http.MethodGet)
	movieRouter.HandleFunc("/{id}", web.AccessGuard(controller.deleteMovieById, true)).Methods(http.MethodDelete)
	movieRouter.HandleFunc("/{id}", web.AccessGuard(controller.updateMovieById, true)).Methods(http.MethodPut)
}

func NewMovieController(service *service.MovieService) *MovieController {
	return &MovieController{
		service: service,
	}
}

func (controller *MovieController) createMovie(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header with the boundary parameter
	contentType := r.Header.Get("Content-Type")
	boundary := web.ExtractBoundary(contentType)

	// Check if the Content-Type is multipart/form-data and if the boundary exists
	if !strings.HasPrefix(contentType, "multipart/form-data") || boundary == "" {
		web.RespondJSON(w, http.StatusBadRequest, "Request Content-Type should be multipart/form-data with boundary parameter")
		return
	}

	// Parse the multipart form with the provided boundary
	err := r.ParseMultipartForm(10 << 20) // Max memory allocated for parsing multipart form data (10 MB)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, "Error parsing multipart form: "+err.Error())
		return
	}

	// Retrieve the movie data from the form
	var movie model.Movie
	if err := web.UnmarshalForm(r.Form, &movie); err != nil {
		web.RespondJSON(w, http.StatusBadRequest, "Error parsing movie data: "+err.Error())
		return
	}

	// Upload image and get upload string
	imageURL, err := cloudinary.CloudinaryUploadImage(r)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, "Error uploading image to cloudinary: "+err.Error())
		return
	}

	token := security.TokenFromContext(r.Context())
	if !token.IsAdmin {
		web.RespondJSON(w, http.StatusUnauthorized, "User unauthorized to access this route")
		return
	}

	// Assign CreatedBy as the user's id from token and imageURL
	movie.CreatedBy = token.ID
	movie.ImageUrl = imageURL

	// Create the movie
	err = controller.service.Create(&movie)
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

func (controller *MovieController) updateMovieById(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	web.UnmarshalJSON(r, &movie)

	vars := mux.Vars(r)

	id := vars["id"]
	movie.ID = uuid.FromStringOrNil(id)

	err := controller.service.UpdateMovie(&movie)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, "Movie Updated Successfully.")
}
