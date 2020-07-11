package checkins

import (
	"context"
	"time"

	t "github.com/contact-tracker/apiService/check-ins/types"
	pT "github.com/contact-tracker/apiService/places/types"
	uT "github.com/contact-tracker/apiService/users/types"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

type repository interface {
	Get(ctx context.Context, id string) (*t.CheckIn, error)
	GetAll(ctx context.Context, userID, placeID *string, start, end *time.Time) ([]*t.CheckIn, error)
	LastCheckIn(ctx context.Context, userID, placeID string) (*t.CheckIn, error)
	Create(ctx context.Context, checkIn *t.CheckIn) (*t.CheckIn, error)
	CheckOut(ctx context.Context, id string) (*t.CheckIn, error)
	Delete(ctx context.Context, id string) error
}

type rpc interface {
	GetPlace(ctx context.Context, id string) (*pT.Place, error)
	GetUser(ctx context.Context, id string) (*uT.User, error)
}

// Usecase for interacting with users
type Usecase struct {
	Repository repository
	RPC        rpc
}

// Get a single check ins
func (u *Usecase) Get(ctx context.Context, id string) (*t.CheckIn, error) {
	checkIn, err := u.Repository.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching a single check in")
	}
	return checkIn, nil
}

// GetAll gets all check ins
func (u *Usecase) GetAll(ctx context.Context, req *t.GetCheckIns) ([]*t.CheckIn, error) {
	validate = validator.New()
	if err := validate.Struct(*req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	checkIns, err := u.Repository.GetAll(ctx, req.UserID, req.PlaceID, req.Start, req.End)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching check ins")
	}
	return checkIns, nil
}

// CheckIn or CheckOut based on user and place ID
func (u *Usecase) CheckIn(ctx context.Context, req *t.CreateCheckIn) (resp *t.CheckIn, err error) {
	validate = validator.New()
	if err := validate.Struct(*req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	var checkIn *t.CheckIn
	if last, err := u.Repository.LastCheckIn(ctx, req.UserID, req.PlaceID); err != nil {
		return nil, err
	} else if last != nil {
		resp, err = u.Repository.CheckOut(ctx, last.ID)
		if err != nil {
			return nil, err
		}
	} else {
		user, err := u.RPC.GetUser(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		place, err := u.RPC.GetPlace(ctx, req.PlaceID)
		if err != nil {
			return nil, err
		}
		checkIn = t.NewCheckIn(t.ToCheckInPlace(place), t.ToCheckInUser(user))
		resp, err = u.Repository.Create(ctx, checkIn)
	}
	return
}

func (u *Usecase) newID() string {
	uid := uuid.New()
	return uid.String()
}
