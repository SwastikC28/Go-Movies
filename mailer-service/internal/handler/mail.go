package handler

import (
	"context"
	"fmt"
	"mailer-service/internal/constants"
	"mailer-service/internal/model"
	"mailer-service/internal/service"
	"mailer-service/internal/utils"
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
		Data:        constants.GetRegisterHTMLBody(user.Name),
		DataMap:     map[string]interface{}{},
	}

	data := map[string]interface{}{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := utils.BuildHTMLMessage(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	plainMessage, err := utils.BuildPlainTextMessage(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = m.MailService.SendMail(msg, formattedMessage, plainMessage)
	if err != nil {
		fmt.Println(err)
	}
}
