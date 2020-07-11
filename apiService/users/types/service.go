package types

import (
	"fmt"

	"github.com/contact-tracker/apiService/pkg/email"
)

// SignInReq - request to check in user at place
type SignInReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func WelcomeEmailInput(user *User, usersHost string) *email.EmailInput {
	confLink := fmt.Sprintf("%s/users/%s/confirm", usersHost, user.ID)
	body := fmt.Sprintf("Welcome to contract tracker %s!\\n\\nA new user account has been created. Please follow this link to confirm your email: %s", user.Name, confLink)

	return &email.EmailInput{
		ToAddresses: []*string{&user.Email},
		HtmlBody:    body,
		TextBody:    body,
		Subject:     "Welcome to contract tracker!",
	}
}
