package checkins

import (
	"context"

	"go.uber.org/zap"

	t "github.com/contact-tracker/apiService/check-ins/types"
)

// LoggerAdapter wraps the usecase interface
// with a logging adapter which can be swapped out
type LoggerAdapter struct {
	Logger  *zap.Logger
	Usecase CheckInService
}

func (a *LoggerAdapter) logErr(err error) {
	if err != nil {
		a.Logger.Error(err.Error())
	}
}

// Get a single checkIn
func (a *LoggerAdapter) Get(ctx context.Context, id string) (*t.CheckIn, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("id", id))
	a.Logger.Info("getting a single checkIn")
	checkIn, err := a.Usecase.Get(ctx, id)
	a.logErr(err)
	return checkIn, err
}

// GetAll gets all checkIns
func (a *LoggerAdapter) GetAll(ctx context.Context, req *t.GetCheckIns) ([]*t.CheckIn, error) {
	defer a.Logger.Sync()
	a.Logger.Info("getting all checkIns")
	checkIns, err := a.Usecase.GetAll(ctx, req)
	a.logErr(err)
	return checkIns, err
}

// CheckIn a user
func (a *LoggerAdapter) CheckIn(ctx context.Context, req *t.CreateCheckIn) (*t.CheckIn, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("userId", req.UserID))
	a.Logger.Info("checkin a single user")
	resp, err := a.Usecase.CheckIn(ctx, req)
	a.logErr(err)
	return resp, err
}
