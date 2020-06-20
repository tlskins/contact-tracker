package types

import (
	"fmt"
	"time"

	"github.com/contact-tracker/apiService/pkg/email"
)

// CheckInReq - request to check in user at place
type CheckInReq struct {
	PlaceID string     `json:"placeId" validate:"required"`
	In      *time.Time `json:"in" validate:"required"`
}

// CheckOutReq - request to check out user at place
type CheckOutReq struct {
	CheckInID string     `json:"checkInId" validate:"required"`
	Out       *time.Time `json:"out" validate:"required"`
}

// SignInReq - request to check in user at place
type SignInReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func WelcomeEmailInput(user *User, usersHost string) *email.EmailInput {
	confLink := fmt.Sprintf("%s/users/%s/confirm", usersHost, user.ID)
	body := fmt.Sprintf("Welcome to contract tracker %s!\n\nPlease follow this link to confirm your email: %s", user.Name, confLink)

	return &email.EmailInput{
		ToAddresses: []*string{&user.Email},
		HtmlBody:    body,
		TextBody:    body,
		Subject:     "Welcome to contract tracker!",
	}
}
