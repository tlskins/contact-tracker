package email

import (
	"crypto/tls"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

type EmailClient struct {
	from     string
	password string
	smtpHost string
	smtpPort int
}

func NewEmailClient(from, password, smtpHost, smtpPort string) (*EmailClient, error) {
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return nil, err
	}
	return &EmailClient{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: port,
	}, nil
}

func (c EmailClient) SendEmail(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", c.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)
	dialer := gomail.NewDialer(c.smtpHost, c.smtpPort, c.from, c.password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return dialer.DialAndSend(msg)
}
