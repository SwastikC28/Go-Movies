package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type PublishEvent struct {
	payload      interface{}
	retryCount   int
	routingKey   string
	corelationID string
}

type RetryPublishEvent struct {
	Event      *PublishEvent
	RetryCount int
}

type RabbitMQDispatcher struct {
	exchangeKey   string
	connection    *amqp.Connection
	channel       *amqp.Channel
	publishChan   chan *PublishEvent
	retryChan     chan *PublishEvent
	maxRetry      int
	retryDuration time.Duration
}

func NewDispatchHandler(exchangeKey string) (*RabbitMQDispatcher, error) {
	dispatcher := &RabbitMQDispatcher{
		exchangeKey:   exchangeKey,
		publishChan:   make(chan *PublishEvent),
		retryChan:     make(chan *PublishEvent),
		maxRetry:      3,
		retryDuration: 5 * time.Second,
	}

	conn, channel := connectRabbitMQ(dispatcher.exchangeKey)
	dispatcher.connection = conn
	dispatcher.channel = channel

	return dispatcher, nil
}

// Start starts the dispatcher goroutine to listen for messages to be published
func (dispatcher *RabbitMQDispatcher) Start(context context.Context) {

	// Listening to PublishChan
	go func() {
		for msg := range dispatcher.publishChan {
			var err error
			body, ok := msg.payload.([]byte)
			if !ok {
				// Marshal Body
				body, err = json.Marshal(msg.payload)
				if err != nil {
					fmt.Println("Error Marshalling your event", err)
				}
			}

			err = dispatcher.channel.Publish(
				dispatcher.exchangeKey, // exchange
				msg.routingKey,         // routing key
				false,                  // mandatory
				false,                  // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				})
			if err != nil {
				fmt.Println("Error Publishing your event", err)
				go dispatcher.retry(msg)
			}

			log.Printf("Published message with routing key '%s'", msg.routingKey)
		}
	}()

	// Listening to RetryChan
	go func() {
		for msg := range dispatcher.retryChan {

			if msg.retryCount > dispatcher.maxRetry {
				continue
			}

			var err error
			body, ok := msg.payload.([]byte)
			if !ok {
				// Marshal Body
				body, err = json.Marshal(msg.payload.([]byte))
				if err != nil {
					fmt.Println("Error Marshalling your event", err)
				}
			}

			err = dispatcher.channel.Publish(
				dispatcher.exchangeKey, // exchange
				msg.routingKey,         // routing key
				false,                  // mandatory
				false,                  // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				})
			if err != nil {
				fmt.Println("Error Publishing your event", err)
				go dispatcher.retry(msg)
			}

			log.Printf("Published message with routing key '%s'", msg.routingKey)
		}
	}()
}

func (dispatcher *RabbitMQDispatcher) Stop(context context.Context) {
	close(dispatcher.publishChan) // Close the publishChan to signal termination
	close(dispatcher.retryChan)   // Close the retryChan to signal termination
	dispatcher.channel.Close()
	dispatcher.connection.Close()
}

func (dispatcher *RabbitMQDispatcher) Publish(corelationId string, topic string, payload interface{}) {
	event := &PublishEvent{
		routingKey:   topic,
		payload:      payload,
		corelationID: corelationId,
	}

	dispatcher.publishChan <- event
}

func (dispatcher *RabbitMQDispatcher) retry(publishEvent *PublishEvent) {
	publishEvent.retryCount++
	time.Sleep(dispatcher.retryDuration)
	dispatcher.retryChan <- publishEvent
}

func connectRabbitMQ(exchangeName string) (*amqp.Connection, *amqp.Channel) {
	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err == nil {
			ch, err := conn.Channel()
			if err != nil {
				fmt.Println("Failed to open a channel", err)
			} else {
				err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
				if err != nil {
					fmt.Println("Failed to declare an exchange", err)
				} else {
					fmt.Println("Successfully Connected to RabbitMQ")
					return conn, ch
				}
			}

			conn.Close()
		}

		fmt.Println("Cannot connect to RabbitMQ. Retrying in 5 seconds ", err)
		time.Sleep(5 * time.Second)
	}
}
