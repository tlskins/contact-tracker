package users

import (
	"context"

	"go.uber.org/zap"

	t "github.com/contact-tracker/apiService/users/types"
)

// LoggerAdapter wraps the usecase interface
// with a logging adapter which can be swapped out
type LoggerAdapter struct {
	Logger  *zap.Logger
	Usecase UserService
}

func (a *LoggerAdapter) logErr(err error) {
	if err != nil {
		a.Logger.Error(err.Error())
	}
}

// Get a single user
func (a *LoggerAdapter) Get(ctx context.Context, id string) (*t.User, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("id", id))
	a.Logger.Info("getting a single user")
	user, err := a.Usecase.Get(ctx, id)
	a.logErr(err)
	return user, err
}

// GetAll gets all users
func (a *LoggerAdapter) GetAll(ctx context.Context) ([]*t.User, error) {
	defer a.Logger.Sync()
	a.Logger.Info("getting all users")
	users, err := a.Usecase.GetAll(ctx)
	a.logErr(err)
	return users, err
}

// Update a single user
func (a *LoggerAdapter) Update(ctx context.Context, user *t.UpdateUser) (*t.User, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("id", user.ID))
	a.Logger.Info("updating a single user")
	resp, err := a.Usecase.Update(ctx, user)
	a.logErr(err)
	return resp, err
}

// CheckIn a single user
func (a *LoggerAdapter) CheckIn(ctx context.Context, id string, req *t.CheckInReq) (*t.User, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("id", id))
	a.Logger.Info("checkin a single user")
	resp, err := a.Usecase.CheckIn(ctx, id, req)
	a.logErr(err)
	return resp, err
}

// CheckOut a single user
func (a *LoggerAdapter) CheckOut(ctx context.Context, id string, req *t.CheckOutReq) (*t.User, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("id", id))
	a.Logger.Info("checkout a single user")
	resp, err := a.Usecase.CheckOut(ctx, id, req)
	a.logErr(err)
	return resp, err
}

// Create a single user
func (a *LoggerAdapter) Create(ctx context.Context, req *t.CreateUser) (*t.User, error) {
	defer a.Logger.Sync()
	a.Logger.Info("creating a single user")
	usr, err := a.Usecase.Create(ctx, req)
	a.logErr(err)
	return usr, err
}

// Delete a single user
func (a *LoggerAdapter) Delete(ctx context.Context, id string) error {
	defer a.Logger.Sync()
	a.Logger.Info("deleting a single user")
	err := a.Usecase.Delete(ctx, id)
	a.logErr(err)
	return err
}

// SignIn a single user
func (a *LoggerAdapter) SignIn(ctx context.Context, req *t.SignInReq) (*t.User, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("email", req.Email))
	a.Logger.Info("sign in a single user")
	resp, err := a.Usecase.SignIn(ctx, req)
	a.logErr(err)
	return resp, err
}

// Confirm a single user
func (a *LoggerAdapter) Confirm(ctx context.Context, id string) error {
	defer a.Logger.Sync()
	a.Logger.Info("confirming a single user")
	err := a.Usecase.Confirm(ctx, id)
	a.logErr(err)
	return err
}
