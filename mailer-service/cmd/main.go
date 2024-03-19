package main

import (
	"fmt"
	"log"
	"math"
	"time"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to RabbitMQ
	amqp, err := connectRabbitMQ()
	if err != nil {
		fmt.Println(err)
	}

	// Close AMQP Connection
	defer func() {
		err := amqp.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Listen for Messages

	// Create Consumer

	// Watch the queue and consume events
}

func connectRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	// dont continue until rabbit is ready
	for {
		// url := fmt.Sprintf("amqp://%s:%s@%s:5672/", os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PASSWORD"))
		url := "amqp://guest@rabbitmq"
		c, err := amqp.Dial(url)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backoff)
		continue

	}

	log.Println("RabbitMQ connected successfully")
	return connection, nil
}
