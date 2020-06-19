package types

import (
	"time"

	"github.com/contact-tracker/apiService/pkg/auth"
)

// User -
type User struct {
	ID                string     `bson:"_id" json:"id"`
	Email             string     `bson:"em" json:"email" validate:"email,required"`
	Name              string     `bson:"nm" json:"name" validate:"required,gte=1,lte=50"`
	EncryptedPassword string     `bson:"pwd" json:"encryptedPassword" validate:"required"`
	Confirmed         bool       `bson:"conf" json:"confirmed"`
	LastLoggedIn      *time.Time `bson:"lstLogIn" json:"lastLoggedIn"`
	CheckIns          []*CheckIn `bson:"chks" json:"checkIns"`
}

// UpdateUser -
type UpdateUser struct {
	ID           string     `bson:"-" json:"id"`
	Email        *string    `bson:"em,omitempty" json:"email,omitempty"`
	Name         *string    `bson:"nm,omitempty" json:"name,omitempty" validate:"gte=1,lte=50"`
	LastLoggedIn *time.Time `bson:"lstLogIn,omitempty" json:"-"`
}

// CreateUser -
type CreateUser struct {
	Email    string `json:"email"`
	Name     string `json:"name" validate:"required,gte=3,lte=50"`
	Password string `json:"password" validate:"gte=3,lte=50"`
}

func (c CreateUser) ToUser(newID string) (*User, error) {
	encPwd, err := auth.EncryptPassword(c.Password)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:                newID,
		Email:             c.Email,
		Name:              c.Name,
		EncryptedPassword: encPwd,
	}, nil
}
