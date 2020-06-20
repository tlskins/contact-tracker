package types

import (
	"time"
)

// User -
type User struct {
	ID                string     `bson:"_id" json:"id"`
	Email             string     `bson:"em" json:"email" validate:"email,required"`
	Name              string     `bson:"nm" json:"name" validate:"required,gte=1,lte=50"`
	EncryptedPassword string     `bson:"pwd" json:"-"`
	Confirmed         bool       `bson:"conf" json:"confirmed"`
	LastLoggedIn      *time.Time `bson:"lstLogIn" json:"lastLoggedIn"`
	CheckIns          []*CheckIn `bson:"chks" json:"checkIns"`
}

func (u User) GetAuthables() (id, email string, conf bool) {
	return u.ID, u.Email, u.Confirmed
}

// UpdateUser -
type UpdateUser struct {
	ID           string     `bson:"-" json:"id"`
	Email        *string    `bson:"em,omitempty" json:"email,omitempty"`
	Name         *string    `bson:"nm,omitempty" json:"name,omitempty" validate:"gte=1,lte=50"`
	LastLoggedIn *time.Time `bson:"lstLogIn,omitempty" json:"-"`
	Confirmed    *bool      `bson:"conf,omitempty" json:"-"`
}

// CreateUser -
type CreateUser struct {
	Email    string `json:"email"`
	Name     string `json:"name" validate:"required,gte=3,lte=50"`
	Password string `json:"password" validate:"gte=3,lte=50"`
}

func (c CreateUser) ToUser() *User {
	return &User{
		Email: c.Email,
		Name:  c.Name,
	}
}
