package places

import (
	"context"

	t "github.com/contact-tracker/apiService/places/types"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

type repository interface {
	Get(ctx context.Context, id string) (*t.Place, error)
	GetAll(ctx context.Context) ([]*t.Place, error)
	Update(ctx context.Context, place *t.UpdatePlace) (*t.Place, error)
	Create(ctx context.Context, place *t.Place) (*t.Place, error)
	Delete(ctx context.Context, id string) error
}

// Usecase for interacting with places
type Usecase struct {
	Repository repository
}

// Get a single place
func (u *Usecase) Get(ctx context.Context, id string) (*t.Place, error) {
	place, err := u.Repository.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching a single place")
	}
	return place, nil
}

// GetAll gets all places
func (u *Usecase) GetAll(ctx context.Context) ([]*t.Place, error) {
	places, err := u.Repository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching all places")
	}
	return places, nil
}

// Update a single place
func (u *Usecase) Update(ctx context.Context, place *t.UpdatePlace) (resp *t.Place, err error) {
	validate = validator.New()
	if err = validate.Struct(place); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	if resp, err = u.Repository.Update(ctx, place); err != nil {
		return nil, errors.Wrap(err, "error updating place")
	}
	return resp, nil
}

// Create a single place
func (u *Usecase) Create(ctx context.Context, place *t.Place) (resp *t.Place, err error) {
	validate = validator.New()
	if err := validate.Struct(*place); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	place.ID = u.newID()
	if resp, err = u.Repository.Create(ctx, place); err != nil {
		return nil, errors.Wrap(err, "error creating new place")
	}

	return resp, nil
}

// Delete a single place
func (u *Usecase) Delete(ctx context.Context, id string) error {
	if err := u.Repository.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "error deleting place")
	}
	return nil
}

func (u *Usecase) newID() string {
	uid := uuid.New()
	return uid.String()
}
