package model

import (
	"shared/pkg/model"
)

type Order struct {
	model.Base
	Entity   string `json:"entity"`
	Amount   string `json:"amount"`
	IsPaid   bool   `json:"isPaid"`
	Currency string `json:"currency"`
	Receipt  string `json:"receipt"`
	OrderId  string `json:"orderId"`
}

type Payment struct {
	model.Base
	OrderId           string `json:"orderId"`
	RentalId          string `json:"rentalId"`
	RazorpayPaymentId string `json:"razorpay_payment_id"`
	RazorpayOrderId   string `json:"razorpay_order_id"`
	RazorpaySignature string `json:"razorpay_signature"`
}
