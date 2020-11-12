package places

import (
	"context"
	"fmt"
	"time"

	"github.com/contact-tracker/apiService/pkg/auth"
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
	FindByEmail(ctx context.Context, email string) (*t.Place, error)
}

// Usecase for interacting with places
type Usecase struct {
	Repository repository
	placesHost string
	storePwd   string
	adminPlace *t.Place
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
func (u *Usecase) Create(ctx context.Context, req *t.CreatePlace) (resp *t.Place, err error) {
	validate = validator.New()
	if err := validate.Struct(*req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	if _, err := u.Repository.FindByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("error place for email already exists")
	}

	place := req.ToPlace()
	place.ID = u.newID()
	now := time.Now()
	place.LastLoggedIn = &now
	if place.EncryptedPassword, err = auth.EncryptPassword(req.Password); err != nil {
		return nil, errors.Wrap(err, "error encrypting password")
	}

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

// SignIn -
func (u *Usecase) SignIn(ctx context.Context, req *t.SignInReq) (place *t.Place, err error) {
	validate = validator.New()
	if err = validate.Struct(req); err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	// place, err := u.Repository.FindByEmail(ctx, req.Email)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "error finding place by email")
	// }
	// if err = auth.ValidateCredentials(place.EncryptedPassword, req.Password); err != nil {
	// 	return nil, errors.Wrap(err, "error validating credentials")
	// }

	// now := time.Now()
	// if resp, err = u.Repository.Update(ctx, &t.UpdatePlace{ID: place.ID, LastLoggedIn: &now}); err != nil {
	// 	return nil, errors.Wrap(err, "error updating place")
	// }

	if req.Password != u.storePwd {
		return nil, errors.New("Invalid store owner password")
	}
	place = u.adminPlace
	return
}

// Confirm a place
func (u *Usecase) Confirm(ctx context.Context, id string) error {
	place, err := u.Repository.Get(ctx, id)
	if err != nil {
		errors.Wrap(err, "error getting place for confirmation")
	}
	if place.Confirmed {
		return fmt.Errorf("place already confirmed")
	}

	c := true
	if _, err := u.Repository.Update(ctx, &t.UpdatePlace{ID: id, Confirmed: &c}); err != nil {
		return errors.Wrap(err, "error confirming place")
	}
	return nil
}

func (u *Usecase) newID() string {
	uid := uuid.New()
	return uid.String()
}
