package types

import (
	"time"

	pT "github.com/contact-tracker/apiService/places/types"

	"github.com/google/uuid"
)

// CheckIn - check in and out
type CheckIn struct {
	ID    string     `bson:"id" json:"id"`
	In    *time.Time `bson:"in" json:"in" validate:"required"`
	Out   *time.Time `bson:"out" json:"out"`
	Place *Place     `bson:"place" json:"place" validate:"required"`
}

func NewCheckIn(place *Place, in *time.Time) *CheckIn {
	return &CheckIn{
		ID:    newCheckInID(),
		In:    in,
		Place: place,
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

func ToUserPlace(place *pT.Place) *Place {
	return &Place{
		ID:   place.ID,
		Name: place.Name,
	}
}
