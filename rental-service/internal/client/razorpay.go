package client

import (
	"fmt"
	"os"
	"time"

	razorpay "github.com/razorpay/razorpay-go"
)

type RazorPayOrder struct {
}

func CreateRazorPayOrder(amount int) (map[string]interface{}, error) {

	razorpayKey := os.Getenv("RAZORPAY_KEY")
	razorpaySecret := os.Getenv("RAZORPAY_SECRET")

	client := razorpay.NewClient(razorpayKey, razorpaySecret)

	uniqueNumber := (time.Now().UnixNano() % 10000000000)

	// Generate random receipt
	receipt := fmt.Sprintf("receipt_%v", uniqueNumber)

	// Here razor pay takes money in paise
	// Multiply by 100

	// Set Order Amount and Receipt
	data := map[string]interface{}{
		"amount":          amount * 100,
		"currency":        "INR",
		"receipt":         receipt,
		"partial_payment": false,
	}

	// Create order from razorpay
	order, err := client.Order.Create(data, nil)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return order, nil
}
