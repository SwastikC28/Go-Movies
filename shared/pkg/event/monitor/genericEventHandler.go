package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var queueRoutes = map[string]func(ctx context.Context, evt *EventInfo){}

type GenericEventHandler struct {
	Queue       string
	ExchangeKey string
	Conn        *amqp.Connection
	Channel     *amqp.Channel
	CloseChan   chan struct{}
}

func NewGenericEventHandler(queueName, exchangeKey string) EventHandler {
	return &GenericEventHandler{
		Queue:       queueName,
		ExchangeKey: exchangeKey,
	}
}

func (e *GenericEventHandler) Start(context.Context) {
	e.connectRabbitMQ()

	// create exchange
	e.createExchange()

	// create queue
	e.createQueue()

	// Listen for event
	e.listen()
}

func (e *GenericEventHandler) Stop(ctx context.Context) {
	select {
	case <-e.CloseChan:
		// Close the channel and RabbitMQ connection
		if e.Channel != nil {
			if err := e.Channel.Close(); err != nil {
				fmt.Println("Error closing channel:", err)
			}
		}
		if e.Conn != nil {
			if err := e.Conn.Close(); err != nil {
				fmt.Println("Error closing connection:", err)
			}
		}
	case <-ctx.Done():
		// If the context is canceled before receiving from CloseChan, perform cleanup here
		fmt.Println("Stop function canceled before receiving from CloseChan")
	}
}

func (e *GenericEventHandler) connectRabbitMQ() {
	for {
		con, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println(err)
			fmt.Println("retrying connecting to rabbitmq. Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		channel, err := con.Channel()
		if err != nil {
			fmt.Println(err)
			fmt.Println("retrying creating a new channel. Retrying in 5 seconds...")
			continue
		}

		e.Conn = con
		e.Channel = channel
		fmt.Println("Connected to RabbitMQ Successfully")
		return
	}
}

func (e *GenericEventHandler) createExchange() {
	for {
		err := e.Channel.ExchangeDeclare(e.ExchangeKey, "topic", true, false, false, false, nil)
		if err != nil {
			fmt.Println("Failed to create exchange. Retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	fmt.Println("Exchange created successfully.")
}

func (e *GenericEventHandler) createQueue() {
	for {
		_, err := e.Channel.QueueDeclare(e.Queue, true, false, false, false, nil)
		if err != nil {
			log.Printf("Failed to create queue: %v. Retrying in 5 seconds", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Bind all events with exchange
		for event := range queueRoutes {
			err = e.Channel.QueueBind(e.Queue, event, e.ExchangeKey, false, nil)
			if err != nil {
				log.Printf("Failed to bind queue with event %s: %v", event, err)
			}
		}
		break
	}
}

func (e *GenericEventHandler) listen() {
	logs, err := e.Channel.Consume(e.Queue, "", true, false, false, false, nil)
	if err != nil {
		fmt.Println("Failed to register a consumer for movies queue", err)
		return
	}

	go func() {
		for {
			select {
			case log := <-logs:
				handler, ok := queueRoutes[log.RoutingKey]

				if !ok {
					fmt.Printf("No handler found for event type: %s\n", log.RoutingKey)
					continue
				}

				evt := EventInfo{}
				err := json.Unmarshal(log.Body, &evt)
				if err != nil {
					fmt.Println("Error Unmarshalling data")
					continue
				}

				handler(context.Background(), &evt)

			case <-e.CloseChan:
				// Close the channel when the signal is received
				if err := e.Channel.Cancel("", false); err != nil {
					fmt.Println("Error canceling consumer:", err)
				}
				return
			}

		}
	}()
}

func AddRoute(evtName string, handler func(ctx context.Context, evt *EventInfo)) {
	queueRoutes[evtName] = handler
}
