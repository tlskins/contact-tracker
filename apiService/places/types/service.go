package types

import (
	"fmt"

	"github.com/contact-tracker/apiService/pkg/email"
)

// SignInReq - request to check in place at place
type SignInReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func WelcomeEmailInput(place *Place, placesHost string) *email.EmailInput {
	confLink := fmt.Sprintf("%s/places/%s/confirm", placesHost, place.ID)
	body := fmt.Sprintf("Welcome to contract tracker %s!\\n\\nA new place has been created. Please follow this link to confirm your email: %s", place.Name, confLink)

	return &email.EmailInput{
		ToAddresses: []*string{&place.Email},
		HtmlBody:    body,
		TextBody:    body,
		Subject:     "Welcome to contract tracker!",
	}
}
