package controller

import (
	"fmt"
	"net/http"
	"rental-service/internal/model"
	"rental-service/internal/service"
	"shared/datastore"
	"shared/middleware"
	"shared/pkg/web"
	"shared/security"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type RentalController struct {
	service *service.RentalService
}

func (controller *RentalController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----Rental Controller Registered-----")
	rentalRouter := router.PathPrefix("/rental").Subrouter()

	rentalRouter.Use(middleware.ReqLogger)

	rentalRouter.HandleFunc("/myrentals", web.AccessGuard(controller.getMyRentals, false)).Methods(http.MethodGet)
	rentalRouter.HandleFunc("/{movieId}/{userId}", web.AccessGuard(controller.createRental, false)).Methods(http.MethodPost)
	rentalRouter.HandleFunc("", web.AccessGuard(controller.getRentals, true)).Methods(http.MethodGet)
	rentalRouter.HandleFunc("/{id}", web.AccessGuard(controller.getRentalById, false)).Methods(http.MethodGet)
	rentalRouter.HandleFunc("/{id}", web.AccessGuard(controller.deleteRentalById, true)).Methods(http.MethodDelete)
}

func NewRentalController(service *service.RentalService) *RentalController {
	return &RentalController{
		service: service,
	}
}

func (controller *RentalController) createRental(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId := params["userId"]
	movieId := params["movieId"]

	rental := model.Rental{}

	rental.Status = "unpaid"
	rental.RentalDate = time.Now()
	rental.DueDate = time.Now().Add(time.Hour * 24 * 7)
	rental.UserId = uuid.FromStringOrNil(userId)
	rental.MovieId = uuid.FromStringOrNil(movieId)

	token := security.TokenFromContext(r.Context())

	if token.ID.String() != userId && !token.IsAdmin {
		web.RespondJSON(w, http.StatusUnauthorized, "user unauthorized to access this route")
		return
	}

	err := controller.service.Create(&rental)
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			web.RespondJSON(w, http.StatusNotFound, err.Error())
			return
		}
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, rental)
}

func (controller *RentalController) getRentals(w http.ResponseWriter, r *http.Request) {
	var rentals []model.Rental

	// Get all associations
	matchedAssociations := web.ParseAssociation(r)

	err := controller.service.GetAllRentals(&rentals, matchedAssociations, nil)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, rentals)
}

func (controller *RentalController) getRentalById(w http.ResponseWriter, r *http.Request) {
	var Rental model.Rental

	vars := mux.Vars(r)

	id := vars["id"]

	queryProcessor := []datastore.QueryProcessor{}
	queryProcessor = append(queryProcessor, datastore.Filter("ID = ?", (id)))

	err := controller.service.GetRental(&Rental, queryProcessor)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, Rental)
}

func (controller *RentalController) getMyRentals(w http.ResponseWriter, r *http.Request) {
	var rentals []model.Rental

	token := security.TokenFromContext(r.Context())

	// Get all associations
	includes := web.ParseAssociation(r)

	queryProcessor := []datastore.QueryProcessor{}
	queryProcessor = append(queryProcessor, datastore.Filter("user_id = ?", (token.ID)))

	err := controller.service.GetAllRentals(&rentals, includes, queryProcessor)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, rentals)
}

func (controller *RentalController) deleteRentalById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	err := controller.service.DeleteRental(id)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, "Rental Deleted Successfully.")
}
