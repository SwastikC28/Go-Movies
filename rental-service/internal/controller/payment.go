package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"rental-service/internal/client"
	"rental-service/internal/model"
	"rental-service/internal/service"
	"shared/middleware"
	"shared/pkg/event/publisher"
	"shared/pkg/web"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type PaymentController struct {
	service    *service.PaymentService
	dispatcher publisher.Dispatcher
}

func (controller *PaymentController) RegisterRoutes(router *mux.Router) {
	fmt.Println("-----Payment Controller Registered-----")
	paymentRouter := router.PathPrefix("/payment").Subrouter()

	paymentRouter.Use(middleware.ReqLogger)

	paymentRouter.HandleFunc("/verify", controller.VerifyPayment).Methods(http.MethodPost)
	paymentRouter.HandleFunc("/{rentalId}", controller.CreateOrder).Methods(http.MethodPost)
}

func NewPaymentController(service *service.PaymentService, dispatcher publisher.Dispatcher) *PaymentController {
	return &PaymentController{
		service:    service,
		dispatcher: dispatcher,
	}
}

func (controller *PaymentController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	rentalId := params["rentalId"]

	var rental model.Rental
	rental.ID = uuid.FromStringOrNil(rentalId)

	err := controller.service.GetRentalFees(&rental)
	if err != nil {
		web.RespondJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	if rental.Status == "paid" {
		web.RespondJSON(w, http.StatusUnauthorized, "rental fees already paid.")
		return
	}

	// Get order details from razorpay
	order, err := client.CreateRazorPayOrder(rental.RentalFee)
	if err != nil {
		web.RespondJSON(w, http.StatusUnauthorized, err)
		return
	}

	web.RespondJSON(w, http.StatusCreated, order)
}

func (controller *PaymentController) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	var payment model.Payment

	web.UnmarshalJSON(r, &payment)

	computeString := fmt.Sprintf("%s|%s", payment.OrderId, payment.RazorpayPaymentId)
	generated_signature := computeHMACSHA256([]byte(computeString), os.Getenv("RAZORPAY_SECRET"))

	if generated_signature != payment.RazorpaySignature {
		web.RespondJSON(w, http.StatusUnauthorized, "Invalid Signature")
	}

	// Save Payment Details in the database
	err := controller.service.SavePayment(payment)
	if err != nil {
		web.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondJSON(w, http.StatusOK, "Payment was successful")
}

func computeHMACSHA256(data []byte, key string) string {
	// Convert the key to bytes
	keyBytes := []byte(key)

	// Create a new HMAC-SHA256 hasher
	hasher := hmac.New(sha256.New, keyBytes)

	// Write the data to the hasher
	hasher.Write(data)

	// Get the computed hash
	hash := hasher.Sum(nil)

	// Encode the hash to a hexadecimal string
	return hex.EncodeToString(hash)
}
