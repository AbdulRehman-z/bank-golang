package mail

import (
	"bytes"
	"net/smtp"
	"strconv"
)

type Mail interface {
	SendEmail(to []string, subject string, body string) error
}

type MailService struct {
	Host     string
	Port     int
	name     string
	email    string
	Password string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) Mail {
	return &MailService{
		Host:     "smtp.gmail.com",
		Port:     587,
		name:     name,
		email:    fromEmailAddress,
		Password: fromEmailPassword,
	}
}

func (m *MailService) SendEmail(to []string, subject string, body string) error {
	receiver := to
	// Sender data.
	from := m.email
	password := m.Password

	// smtp server configuration.
	smtpHost := m.Host
	smtpPort := m.Port

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	var bodyMessage bytes.Buffer
	bodyMessage.Write([]byte("Subject: " + subject + "\n" + body))

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), // Convert Int to String
		auth,
		from,
		receiver,
		bodyMessage.Bytes())

	if err != nil {
		return err
	}

	return nil

}