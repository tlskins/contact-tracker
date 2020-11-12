package types

import "time"

// Place -
type Place struct {
	ID                string     `bson:"_id" json:"id"`
	Name              string     `bson:"nm" json:"name" validate:"required,gte=1,lte=50"`
	Email             string     `bson:"em" json:"email" validate:"email,required"`
	EncryptedPassword string     `bson:"pwd" json:"-"`
	Confirmed         bool       `bson:"conf" json:"confirmed"`
	LastLoggedIn      *time.Time `bson:"lstLogIn" json:"lastLoggedIn"`
	AuthToken         string     `bson:"-" json:"authToken"`
}

func (u Place) GetAuthables() (id, email string, conf bool) {
	return u.ID, u.Email, u.Confirmed
}

// UpdatePlace -
type UpdatePlace struct {
	ID                string     `bson:"-" json:"id"`
	Name              *string    `bson:"nm,omitempty" json:"name,omitempty" validate:"gte=1,lte=50"`
	Email             *string    `bson:"em,omitempty" json:"email,omitempty"`
	Confirmed         *bool      `bson:"conf,omitempty" json:"-"`
	EncryptedPassword *string    `bson:"pwd,omitempty" json:"-"`
	LastLoggedIn      *time.Time `bson:"lstLogIn,omitempty" json:"-"`
}

// CreatePlace -
type CreatePlace struct {
	Email    string `json:"email"`
	Name     string `json:"name" validate:"required,gte=3,lte=50"`
	Password string `json:"password" validate:"gte=3,lte=50"`
}

func (c CreatePlace) ToPlace() *Place {
	return &Place{
		Email: c.Email,
		Name:  c.Name,
	}
}
