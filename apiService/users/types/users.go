package types

// User -
type User struct {
	ID    string `bson:"_id" json:"id"`
	Email string `bson:"em" json:"email" validate:"email,required"`
	Name  string `bson:"nm" json:"name" validate:"required,gte=1,lte=50"`
}

// UpdateUser -
type UpdateUser struct {
	ID    string  `bson:"-" json:"id"`
	Email *string `bson:"em,omitempty" json:"email,omitempty"`
	Name  *string `bson:"nm,omitempty" json:"name,omitempty" validate:"gte=1,lte=50"`
}
