package service

import (
	"fmt"
	"log"
	"mailer-service/internal/model"

	mail "github.com/xhit/go-simple-mail/v2"
)

const htmlBody = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Hello Gophers!</title>
	</head>
	<body>
		<p>Hello User</b>.</p>
		<p><img src="cid:Gopher.png" alt="Go gopher" /></p>
		<p>Image created by Renee French</p>
	</body>
</html>`

type MailService struct {
	mailserver *mail.SMTPServer
}

func NewMailService(mailserver *mail.SMTPServer) *MailService {
	return &MailService{
		mailserver: mailserver,
	}
}

func (mailService *MailService) SendMail(msg model.Message) error {
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

	email.SetBody(mail.TextHTML, htmlBody)

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

	// mailService.mailserver.Username = "test@example.com"
	// mailService.mailserver.Password = "examplepass"
	// mailService.mailserver.Encryption = mail.EncryptionSTARTTLS

	// // Variable to keep alive connection
	// mailService.mailserver.KeepAlive = false

	// // Timeout for connect to SMTP Server
	// mailService.mailserver.ConnectTimeout = 10 * time.Second

	// // Timeout for send the data and wait respond
	// mailService.mailserver.SendTimeout = 10 * time.Second

	// // Set TLSConfig to provide custom TLS configuration. For example,
	// // to skip TLS verification (useful for testing):
	// mailService.mailserver.TLSConfig = nil
	// // &tls.Config{InsecureSkipVerify: true}
}
