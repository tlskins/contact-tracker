package email

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const charSet = "UTF-8"

type EmailService interface {
	SendEmail(input *EmailInput) error
	UsersHost() string
}

type SESClient struct {
	ses         *ses.SES
	senderEmail string
	usersHost   string
}

func (s SESClient) SendEmail(input *EmailInput) error {
	if len(input.Sender) == 0 {
		input.Sender = s.senderEmail
	}
	if _, err := s.ses.SendEmail(input.ToSES()); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return fmt.Errorf("Code: %s Error: %s", aerr.Code(), aerr.Error())
		}
		return err
	}
	return nil
}

func (s SESClient) UsersHost() string {
	return s.usersHost
}

func Init(region, accessKey, accessSecret, senderEmail, usersHost string) (EmailService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			accessSecret,
			"",
		),
	})
	if err != nil {
		return nil, err
	}
	s := ses.New(sess)

	return SESClient{s, senderEmail, usersHost}, nil
}

type EmailInput struct {
	CcAddresses []*string `json:"cc_addresses"`
	ToAddresses []*string `json:"to_addresses"`
	HtmlBody    string    `json:"html_body"`
	TextBody    string    `json:"text_body"`
	Subject     string    `json:"subject"`
	Sender      string    `json:"sender"`
}

func (e EmailInput) ToSES() *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: e.CcAddresses,
			ToAddresses: e.ToAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(e.HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(e.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(e.Subject),
			},
		},
		Source: aws.String(e.Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}
}
