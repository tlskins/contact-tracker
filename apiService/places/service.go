package places

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	m "github.com/contact-tracker/apiService/pkg/mongo"
	r "github.com/contact-tracker/apiService/places/repository"
	t "github.com/contact-tracker/apiService/places/types"
	"github.com/joho/godotenv"

	"go.uber.org/zap"
)

// PlaceService - is the top level signature of this service
type PlaceService interface {
	Get(ctx context.Context, id string) (*t.Place, error)
	GetAll(ctx context.Context) ([]*t.Place, error)
	Update(ctx context.Context, place *t.UpdatePlace) (*t.Place, error)
	Create(ctx context.Context, place *t.Place) (*t.Place, error)
	Delete(ctx context.Context, id string) error
}

// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func InitMongoService() (PlaceService, error) {
	fmt.Println("Init Places Mongo Service...")
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPlace := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	fmt.Printf("Configs: %s %s %s %s\n\n", mongoDBName, mongoHost, mongoPlace, mongoPwd)

	mc, err := m.NewClient(mongoHost, mongoPlace, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}
	repo := r.NewMongoPlaceRepository(mc, mongoDBName)

	logger, _ := zap.NewProduction()

	usecase := &LoggerAdapter{
		Logger:  logger,
		Usecase: &Usecase{Repository: repo},
	}
	return usecase, nil
}
