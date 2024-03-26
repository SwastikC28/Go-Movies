package handler

import (
	"context"
	"fmt"
	"mailer-service/internal/model"
	"mailer-service/internal/service"
	"os"
	"shared/pkg/event/monitor"
)

type MailEventHandler struct {
	MailService *service.MailService
}

func NewMailEventHandler(mailService *service.MailService) monitor.EventRouteProvider {
	return &MailEventHandler{
		MailService: mailService,
	}
}

func (m *MailEventHandler) RegisterEvents() {
	monitor.AddRoute("mailer.register", m.SendRegisterMail)
}

func (m *MailEventHandler) SendRegisterMail(ctx context.Context, evt *monitor.EventInfo) {

	jsonBody := evt.Payload.(map[string]interface{})

	user := model.User{
		Email: jsonBody["email"].(string),
		Name:  jsonBody["name"].(string),
		ID:    jsonBody["id"].(string),
	}

	msg := model.Message{
		From:        os.Getenv("FROM_ADDRESS"),
		FromName:    os.Getenv("FROM_NAME"),
		To:          user.Email,
		Subject:     os.Getenv("REGISTER_MAIL_SUBJECT"),
		Attachments: []string{},
		Data:        nil,
		DataMap:     map[string]interface{}{},
	}

	err := m.MailService.SendMail(msg)
	if err != nil {
		fmt.Println(err)
	}
}
