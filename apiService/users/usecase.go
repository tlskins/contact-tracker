package users

import (
	"context"
	"fmt"
	"time"

	"github.com/contact-tracker/apiService/pkg/auth"
	"github.com/contact-tracker/apiService/pkg/email"
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
	GetByIds(ctx context.Context, id []string) ([]*t.User, error)
	Search(ctx context.Context, search string) ([]*t.User, error)
	GetAll(ctx context.Context) ([]*t.User, error)
	Update(ctx context.Context, user *t.UpdateUser) (*t.User, error)
	FindByEmail(ctx context.Context, email string) (*t.User, error)
	Create(ctx context.Context, user *t.User) (*t.User, error)
	Delete(ctx context.Context, id string) error
}

// Usecase for interacting with users
type Usecase struct {
	Repository  repository
	EmailClient *email.EmailClient
	usersHost   string
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

// Search for users
func (u *Usecase) Search(ctx context.Context, search string) ([]*t.User, error) {
	users, err := u.Repository.Search(ctx, search)
	if err != nil {
		return nil, errors.Wrap(err, "error searching users")
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

// SignIn -
func (u *Usecase) SignIn(ctx context.Context, req *t.SignInReq) (resp *t.User, err error) {
	validate = validator.New()
	if err = validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	user, err := u.Repository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.Wrap(err, "error finding user by email")
	}
	if err = auth.ValidateCredentials(user.EncryptedPassword, req.Password); err != nil {
		return nil, errors.Wrap(err, "error validating credentials")
	}

	now := time.Now()
	if resp, err = u.Repository.Update(ctx, &t.UpdateUser{ID: user.ID, LastLoggedIn: &now}); err != nil {
		return nil, errors.Wrap(err, "error updating user")
	}

	return resp, err
}

// Create a single user
func (u *Usecase) Create(ctx context.Context, req *t.CreateUser) (resp *t.User, err error) {
	validate = validator.New()
	if err := validate.Struct(*req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	if _, err := u.Repository.FindByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("error user for email already exists")
	}

	user := req.ToUser()
	user.ID = u.newID()
	now := time.Now()
	user.LastLoggedIn = &now
	if user.EncryptedPassword, err = auth.EncryptPassword(req.Password); err != nil {
		return nil, errors.Wrap(err, "error encrypting password")
	}

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

// Confirm a single user
func (u *Usecase) Confirm(ctx context.Context, id string) error {
	user, err := u.Repository.Get(ctx, id)
	if err != nil {
		errors.Wrap(err, "error getting user for confirmation")
	}
	if user.Confirmed {
		return fmt.Errorf("user already confirmed")
	}

	c := true
	if _, err := u.Repository.Update(ctx, &t.UpdateUser{ID: id, Confirmed: &c}); err != nil {
		return errors.Wrap(err, "error confirming user")
	}
	return nil
}

// Alert users of possible unsafe contact
func (u *Usecase) AlertUsers(ctx context.Context, ids []string) (err error) {
	var users []*t.User
	if users, err = u.Repository.GetByIds(ctx, ids); err != nil {
		return errors.Wrap(err, "error getting users for alert")
	}
	for _, user := range users {
		if err = u.EmailClient.SendEmail(
			user.Email,
			"Important Notice From Contact Tracker",
			"Hello from Contact Tracker,\n\nThis email is a warning that you may have been in contact with someone who has recently tested positive for COVID-19. Please have yourself tested and be sure to wear a mask in public.\n\nThank You,\nContact Tracker Team",
		); err != nil {
			return
		}
	}
	return nil
}

func (u *Usecase) newID() string {
	uid := uuid.New()
	return uid.String()
}
