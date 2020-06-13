package places

// Place -
type Place struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"nm" json:"name" validate:"required,gte=1,lte=50"`
}

// UpdatePlace -
type UpdatePlace struct {
	ID   string  `bson:"-" json:"id"`
	Name *string `bson:"nm,omitempty" json:"name,omitempty" validate:"gte=1,lte=50"`
}
