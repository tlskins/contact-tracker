package types

import (
	"time"
)

type GetCheckIns struct {
	UserID *string    `bson:"userId" json:"userId"`
	Start  *time.Time `bson:"start" json:"start"`
	End    *time.Time `bson:"end" json:"end"`
}

type CreateCheckIn struct {
	UserID string `bson:"userId" json:"userId" validate:"required"`
}

type CheckInHistory struct {
	ID                string     `bson:"_id" json:"id"`
	In                *time.Time `bson:"in" json:"in"`
	Out               *time.Time `bson:"out" json:"out"`
	User              *User      `bson:"user" json:"user"`
	Contacts          []*CheckIn `bson:"contacts" json:"contacts"`
	TentativeCheckout bool       `bson:"tentative,omitempty" json:"tentativeCheckout,omitempty"`
}

// CheckIn - check in and out
type CheckIn struct {
	ID                string     `bson:"_id" json:"id"`
	In                *time.Time `bson:"in" json:"in" validate:"required"`
	Out               *time.Time `bson:"out" json:"out"`
	User              *User      `bson:"user" json:"user" validate:"required"`
	TentativeCheckout bool       `bson:"tentative,omitempty" json:"tentativeCheckout,omitempty"`
}

// Place - place from user view
type Place struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"nm" json:"name"`
}

// User - user from user view
type User struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"nm" json:"name"`
}
