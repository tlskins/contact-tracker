package types

import (
	"time"

	pT "github.com/contact-tracker/apiService/places/types"
	uT "github.com/contact-tracker/apiService/users/types"

	"github.com/google/uuid"
)

type GetCheckIns struct {
	UserID  *string    `bson:"userId" json:"userId"`
	PlaceID *string    `bson:"placeId" json:"placeId"`
	Start   *time.Time `bson:"start" json:"start"`
	End     *time.Time `bson:"end" json:"end"`
}

type CreateCheckIn struct {
	UserID  string `bson:"userId" json:"userId" validate:"required"`
	PlaceID string `bson:"placeId" json:"placeId" validate:"required"`
}

type CheckInHistory struct {
	ID       string     `bson:"id" json:"id"`
	In       *time.Time `bson:"in" json:"in"`
	Out      *time.Time `bson:"out" json:"out"`
	Place    *Place     `bson:"place" json:"place"`
	User     *User      `bson:"user" json:"user"`
	Contacts []*CheckIn `bson:"contacts" json:"contacts"`
}

// CheckIn - check in and out
type CheckIn struct {
	ID    string     `bson:"id" json:"id"`
	In    *time.Time `bson:"in" json:"in" validate:"required"`
	Out   *time.Time `bson:"out" json:"out"`
	Place *Place     `bson:"place" json:"place" validate:"required"`
	User  *User      `bson:"user" json:"user" validate:"required"`
}

func NewCheckIn(place *Place, user *User) *CheckIn {
	now := time.Now()
	return &CheckIn{
		ID:    newCheckInID(),
		In:    &now,
		Place: place,
		User:  user,
	}
}

func newCheckInID() string {
	uid := uuid.New()
	return uid.String()
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

func ToCheckInPlace(place *pT.Place) *Place {
	return &Place{
		ID:   place.ID,
		Name: place.Name,
	}
}

func ToCheckInUser(user *uT.User) *User {
	return &User{
		ID:   user.ID,
		Name: user.Name,
	}
}
