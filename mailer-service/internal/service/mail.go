package service

import (
	"fmt"
	"log"
	"mailer-service/internal/model"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type MailService struct {
	mailserver *mail.SMTPServer
}

func NewMailService(mailserver *mail.SMTPServer) *MailService {
	return &MailService{
		mailserver: mailserver,
	}
}

func (mailService *MailService) SendMail(msg model.Message, body ...string) error {
	// Configure Mail Client
	mailService.configureMailClient()

	smtpClient, err := mailService.mailserver.Connect()
	if err != nil {
		fmt.Println("Error Connecting to Client", err)
		return err
	}

	// New email simple html with inline and CC
	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject).
		SetListUnsubscribe("<mailto:unsubscribe@example.com?subject=https://example.com/unsubscribe>")

	email.SetBody(mail.TextPlain, body[1])
	email.AddAlternative(mail.TextHTML, body[0])

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Email Sent Successfully")
	return nil
}

func (mailService *MailService) configureMailClient() {
	// SMTP Server
	mailService.mailserver.Host = "mailcatcher"
	mailService.mailserver.Port = 1025
	mailService.mailserver.Encryption = mail.EncryptionNone
	mailService.mailserver.KeepAlive = false
	mailService.mailserver.ConnectTimeout = 10 * time.Second
	mailService.mailserver.SendTimeout = 10 * time.Second

	// mailService.mailserver.Username = "test@example.com"
	// mailService.mailserver.Password = "examplepass"
	// mailService.mailserver.Encryption = mail.EncryptionSTARTTLS

	// // Set TLSConfig to provide custom TLS configuration. For example,
	// // to skip TLS verification (useful for testing):
	// mailService.mailserver.TLSConfig = nil
	// // &tls.Config{InsecureSkipVerify: true}
}
