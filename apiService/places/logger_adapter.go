package places

import (
	"context"

	"go.uber.org/zap"

	t "github.com/contact-tracker/apiService/places/types"
)

// LoggerAdapter wraps the usecase interface
// with a logging adapter which can be swapped out
type LoggerAdapter struct {
	Logger  *zap.Logger
	Usecase PlaceService
}

func (a *LoggerAdapter) logErr(err error) {
	if err != nil {
		a.Logger.Error(err.Error())
	}
}

// Get a single place
func (a *LoggerAdapter) Get(ctx context.Context, id string) (*t.Place, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("id", id))
	a.Logger.Info("getting a single place")
	place, err := a.Usecase.Get(ctx, id)
	a.logErr(err)
	return place, err
}

// GetAll gets all places
func (a *LoggerAdapter) GetAll(ctx context.Context) ([]*t.Place, error) {
	defer a.Logger.Sync()
	a.Logger.Info("getting all places")
	places, err := a.Usecase.GetAll(ctx)
	a.logErr(err)
	return places, err
}

// Update a single place
func (a *LoggerAdapter) Update(ctx context.Context, place *t.UpdatePlace) (*t.Place, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("id", place.ID))
	a.Logger.Info("updating a single place")
	resp, err := a.Usecase.Update(ctx, place)
	a.logErr(err)
	return resp, err
}

// Create a single place
func (a *LoggerAdapter) Create(ctx context.Context, place *t.Place) (*t.Place, error) {
	defer a.Logger.Sync()
	a.Logger.Info("creating a single place")
	usr, err := a.Usecase.Create(ctx, place)
	a.logErr(err)
	return usr, err
}

// Delete a single place
func (a *LoggerAdapter) Delete(ctx context.Context, id string) error {
	defer a.Logger.Sync()
	a.Logger.Info("deleting a single place")
	err := a.Usecase.Delete(ctx, id)
	a.logErr(err)
	return err
}
