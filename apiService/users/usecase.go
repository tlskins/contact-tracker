package users

import (
	"context"

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
	Create(ctx context.Context, user *t.User) (*t.User, error)
	Delete(ctx context.Context, id string) error
}

// Usecase for interacting with users
type Usecase struct {
	Repository repository
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
