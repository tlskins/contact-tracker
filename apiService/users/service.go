package users

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	m "github.com/contact-tracker/apiService/pkg/mongo"
	r "github.com/contact-tracker/apiService/users/repository"
	t "github.com/contact-tracker/apiService/users/types"
	"github.com/joho/godotenv"

	"go.uber.org/zap"
)

// UserService - is the top level signature of this service
type UserService interface {
	Get(ctx context.Context, id string) (*t.User, error)
	GetAll(ctx context.Context) ([]*t.User, error)
	Update(ctx context.Context, user *t.UpdateUser) (*t.User, error)
	Create(ctx context.Context, user *t.User) (*t.User, error)
	Delete(ctx context.Context, id string) error
}

// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func InitMongoService() (UserService, error) {
	fmt.Println("Init Mongo Service...")
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	fmt.Printf("Configs: %s %s %s %s\n\n", mongoDBName, mongoHost, mongoUser, mongoPwd)

	mc, err := m.NewClient(mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}
	repo := r.NewMongoUserRepository(mc, mongoDBName)

	logger, _ := zap.NewProduction()

	usecase := &LoggerAdapter{
		Logger:  logger,
		Usecase: &Usecase{Repository: repo},
	}
	return usecase, nil
}
