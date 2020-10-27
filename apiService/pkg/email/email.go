package email

import (
	"net/smtp"
)

type EmailClient struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewEmailClient(from, password, smtpHost, smtpPort string) *EmailClient {
	return &EmailClient{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (c EmailClient) SendEmail(to, msg string) error {
	msgBytes := []byte(msg)
	recipients := []string{to}
	auth := smtp.PlainAuth("", c.from, c.password, c.smtpHost)
	return smtp.SendMail(c.smtpHost+":"+c.smtpPort, auth, c.from, recipients, msgBytes)
}
