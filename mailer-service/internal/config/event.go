package config

import (
	"context"
	"mailer-service/internal/handler"
	"mailer-service/internal/service"
	"os"
	"shared/pkg/event/monitor"

	_ "github.com/go-sql-driver/mysql"
	mail "github.com/xhit/go-simple-mail/v2"
)

var eventHandler monitor.EventHandler

func StartEventHandler() {
	exchangeKey := os.Getenv("RABBITMQ_EXCHANGE_KEY")

	// Generate Event Handler
	eventHandler = monitor.NewGenericEventHandler("mailer", exchangeKey)

	// Register all the events
	RegisterEvents()

	// Start Event Handler
	eventHandler.Start(context.Background())
}

func StopEventHandler() {
	eventHandler.Stop(context.Background())
}

func RegisterEvents() {
	// Get Mail Service
	mailservice := service.NewMailService(mail.NewSMTPClient())

	// Event Route Provider
	eventRouteProvider := handler.NewMailEventHandler(mailservice)

	// Register Events
	eventRouteProvider.RegisterEvents()
}
