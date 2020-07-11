package checkins

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	repo "github.com/contact-tracker/apiService/check-ins/repository"
	chkRpc "github.com/contact-tracker/apiService/check-ins/rpc"
	t "github.com/contact-tracker/apiService/check-ins/types"
	"github.com/contact-tracker/apiService/pkg/auth"
	m "github.com/contact-tracker/apiService/pkg/mongo"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// CheckInService - is the top level signature of this service
type CheckInService interface {
	Get(ctx context.Context, id string) (*t.CheckIn, error)
	GetAll(ctx context.Context, req *t.GetCheckIns) ([]*t.CheckIn, error)
	CheckIn(ctx context.Context, req *t.CreateCheckIn) (resp *t.CheckIn, err error)
}

// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func Init() (CheckInService, *auth.JWTService, error) {
	fmt.Println("Init CheckIns Mongo Service...")
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoCheckIn := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	usersHost := os.Getenv("USERS_HOST")
	placesHost := os.Getenv("PLACES_HOST")
	jwtKeyPath := os.Getenv("JWT_KEY_PATH")
	jwtSecretPath := os.Getenv("JWT_SECRET_PATH")
	rpcPwd := os.Getenv("RPC_AUTH_PWD")

	// Init mongo repo
	mc, err := m.NewClient(mongoHost, mongoCheckIn, mongoPwd)
	if err != nil {
		log.Fatalf("Error starting mongo client: Error: %v\n", err)
	}
	repository := repo.NewMongoCheckInRepository(mc, mongoDBName)

	// Init rpc client
	rpcClient := chkRpc.NewRPCClient(placesHost, usersHost, rpcPwd)

	// Init jwt service
	jwtKey, err := ioutil.ReadFile(jwtKeyPath)
	if err != nil {
		log.Fatalf("Error reading jwt key file path: %s Error: %v\n", jwtKeyPath, err)
	}
	jwtSecret, err := ioutil.ReadFile(jwtSecretPath)
	if err != nil {
		log.Fatalf("Error reading jwt secret file path: %s Error: %v\n", jwtSecretPath, err)
	}
	j, err := auth.NewJWTService(auth.JWTServiceConfig{
		Key:    jwtKey,
		Secret: jwtSecret,
		RPCPwd: rpcPwd,
	})
	if err != nil {
		log.Fatalf("Error creating jwt service: %v\n", err)
	}

	// Init logger
	logger, _ := zap.NewProduction()

	usecase := &LoggerAdapter{
		Logger: logger,
		Usecase: &Usecase{
			Repository: repository,
			RPC:        rpcClient,
		},
	}
	return usecase, j, nil
}
