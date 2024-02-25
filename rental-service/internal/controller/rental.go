package controller

import (
	"fmt"
	"net/http"
	"rental-service/internal/model"
	"rental-service/internal/service"
	"shared/datastore"
	"shared/middleware"
	"shared/pkg/web"

	"github.com/gorilla/mux"
)

type RentalController struct {
	service *service.RentalService
}

func (controller *RentalController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----Rental Controller Registered-----")
	rentalRouter := router.PathPrefix("/Rental").Subrouter()

	rentalRouter.Use(middleware.ReqLogger)
	rentalRouter.HandleFunc("", controller.createRental).Methods(http.MethodPost)
	rentalRouter.HandleFunc("", controller.getRentals).Methods(http.MethodGet)
	rentalRouter.HandleFunc("/{id}", controller.getRentalById).Methods(http.MethodGet)
	rentalRouter.HandleFunc("/{id}", controller.deleteRentalById).Methods(http.MethodDelete)
}

func NewRentalController(service *service.RentalService) *RentalController {
	return &RentalController{
		service: service,
	}
}

func (controller *RentalController) createRental(w http.ResponseWriter, r *http.Request) {
	var Rental model.Rental
	web.UnmarshalJSON(r, &Rental)

	err := controller.service.Create(&Rental)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, Rental)
}

func (controller *RentalController) getRentals(w http.ResponseWriter, r *http.Request) {
	var Rentals []model.Rental

	err := controller.service.GetAllRentals(&Rentals)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, Rentals)
}

func (controller *RentalController) getRentalById(w http.ResponseWriter, r *http.Request) {
	var Rental model.Rental

	vars := mux.Vars(r)

	id := vars["id"]
	fmt.Println("ID", id)

	queryProcessor := []datastore.QueryProcessor{}
	queryProcessor = append(queryProcessor, datastore.Filter("ID = ?", (id)))

	err := controller.service.GetRental(&Rental, queryProcessor)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, Rental)
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
