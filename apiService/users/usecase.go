package users

import (
	"context"
	"time"

	pT "github.com/contact-tracker/apiService/places/types"
	t "github.com/contact-tracker/apiService/users/types"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

type repository interface {
	Get(ctx context.Context, id string) (*t.User, error)
	GetAll(ctx context.Context) ([]*t.User, error)
	Update(ctx context.Context, user *t.UpdateUser) (*t.User, error)
	CheckIn(ctx context.Context, id string, chk *t.CheckIn) (*t.User, error)
	CheckOut(ctx context.Context, id, chkID string, out *time.Time) (*t.User, error)
	Create(ctx context.Context, user *t.User) (*t.User, error)
	Delete(ctx context.Context, id string) error
}

type rpc interface {
	GetPlace(ctx context.Context, id string) (*pT.Place, error)
}

// Usecase for interacting with users
type Usecase struct {
	Repository repository
	RPC        rpc
}

// Get a single user
func (u *Usecase) Get(ctx context.Context, id string) (*t.User, error) {
	user, err := u.Repository.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching a single user")
	}
	return user, nil
}

// GetAll gets all users
func (u *Usecase) GetAll(ctx context.Context) ([]*t.User, error) {
	users, err := u.Repository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching all users")
	}
	return users, nil
}

// Update a single user
func (u *Usecase) Update(ctx context.Context, user *t.UpdateUser) (resp *t.User, err error) {
	validate = validator.New()
	if err = validate.Struct(user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	if resp, err = u.Repository.Update(ctx, user); err != nil {
		return nil, errors.Wrap(err, "error updating user")
	}
	return resp, nil
}

// CheckIn a single user
func (u *Usecase) CheckIn(ctx context.Context, id string, req *t.CheckInReq) (resp *t.User, err error) {
	validate = validator.New()
	if err = validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	p, err := u.RPC.GetPlace(ctx, req.PlaceID)
	if err != nil {
		return nil, err
	}
	place := t.ToUserPlace(p)

	chk := t.NewCheckIn(place, req.In)

	if resp, err = u.Repository.CheckIn(ctx, id, chk); err != nil {
		return nil, errors.Wrap(err, "error checking in user")
	}
	return resp, nil
}

// CheckOut a single user
func (u *Usecase) CheckOut(ctx context.Context, id string, req *t.CheckOutReq) (resp *t.User, err error) {
	validate = validator.New()
	if err = validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	if resp, err = u.Repository.CheckOut(ctx, id, req.CheckInID, req.Out); err != nil {
		return nil, errors.Wrap(err, "error checking out user")
	}
	return resp, nil
}

// Create a single user
func (u *Usecase) Create(ctx context.Context, user *t.User) (resp *t.User, err error) {
	validate = validator.New()
	if err := validate.Struct(*user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	user.ID = u.newID()
	if resp, err = u.Repository.Create(ctx, user); err != nil {
		return nil, errors.Wrap(err, "error creating new user")
	}

	return resp, nil
}

// Delete a single user
func (u *Usecase) Delete(ctx context.Context, id string) error {
	if err := u.Repository.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "error deleting user")
	}
	return nil
}

func (u *Usecase) newID() string {
	uid := uuid.New()
	return uid.String()
}
