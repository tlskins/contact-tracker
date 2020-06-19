package users

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	m "github.com/contact-tracker/apiService/pkg/mongo"
	repo "github.com/contact-tracker/apiService/users/repository"
	usrRpc "github.com/contact-tracker/apiService/users/rpc"
	t "github.com/contact-tracker/apiService/users/types"
	"github.com/joho/godotenv"

	"go.uber.org/zap"
)

// UserService - is the top level signature of this service
type UserService interface {
	Get(ctx context.Context, id string) (*t.User, error)
	GetAll(ctx context.Context) ([]*t.User, error)
	Update(ctx context.Context, user *t.UpdateUser) (*t.User, error)
	SignIn(ctx context.Context, req *t.SignInReq) (*t.User, error)
	CheckIn(ctx context.Context, id string, chk *t.CheckInReq) (*t.User, error)
	CheckOut(ctx context.Context, id string, req *t.CheckOutReq) (*t.User, error)
	Create(ctx context.Context, user *t.CreateUser) (*t.User, error)
	Delete(ctx context.Context, id string) error
}

// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func Init() (UserService, error) {
	fmt.Println("Init Users Mongo Service...")
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	placesHost := os.Getenv("PLACES_HOST")
	fmt.Printf("Configs: %s %s %s %s %s\n\n", mongoDBName, mongoHost, mongoUser, mongoPwd, placesHost)

	repository, err := initMongoRepo(mongoDBName, mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}

	rpcClient := usrRpc.NewHTTPRPCClient(placesHost)

	logger, _ := zap.NewProduction()

	usecase := &LoggerAdapter{
		Logger: logger,
		Usecase: &Usecase{
			Repository: repository,
			RPC:        rpcClient,
		},
	}
	return usecase, nil
}

func initMongoRepo(mongoDBName, mongoHost, mongoUser, mongoPwd string) (*repo.MongoUserRepository, error) {
	mc, err := m.NewClient(mongoHost, mongoUser, mongoPwd)
	if err != nil {
		return nil, err
	}
	return repo.NewMongoUserRepository(mc, mongoDBName), nil
}
