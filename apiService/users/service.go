package users

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/contact-tracker/apiService/pkg/auth"
	"github.com/contact-tracker/apiService/pkg/email"
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
	Create(ctx context.Context, user *t.CreateUser) (*t.User, error)
	Delete(ctx context.Context, id string) error

	SignIn(ctx context.Context, req *t.SignInReq) (*t.User, error)
	Confirm(ctx context.Context, id string) error
	CheckIn(ctx context.Context, id string, chk *t.CheckInReq) (*t.User, error)
	CheckOut(ctx context.Context, id string, req *t.CheckOutReq) (*t.User, error)
}

// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func Init() (UserService, *auth.JWTService, error) {
	fmt.Println("Init Users Mongo Service...")
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	usersHost := os.Getenv("USERS_HOST")
	placesHost := os.Getenv("PLACES_HOST")
	jwtKeyPath := os.Getenv("JWT_KEY_PATH")
	jwtSecretPath := os.Getenv("JWT_SECRET_PATH")
	sesAccessKey := os.Getenv("AWS_SES_ACCESS_KEY")
	sesAccessSecret := os.Getenv("AWS_SES_ACCESS_SECRET")
	sesRegion := os.Getenv("AWS_SES_REGION")
	senderEmail := os.Getenv("SENDER_EMAIL")
	rpcPwd := os.Getenv("RPC_AUTH_PWD")

	// Init mongo repo
	mc, err := m.NewClient(mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatalf("Error starting mongo client: Error: %v\n", err)
	}
	repository := repo.NewMongoUserRepository(mc, mongoDBName)

	// Init rpc client
	rpcClient := usrRpc.NewRPCClient(placesHost, rpcPwd)

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

	// Init email
	emailClient, err := email.Init(sesRegion, sesAccessKey, sesAccessSecret, senderEmail)
	if err != nil {
		log.Fatalf("Error starting email svc: %v\n", err)
	}

	// Init logger
	logger, _ := zap.NewProduction()

	usecase := &LoggerAdapter{
		Logger: logger,
		Usecase: &Usecase{
			Repository: repository,
			RPC:        rpcClient,
			Email:      emailClient,
			usersHost:  usersHost,
		},
	}
	return usecase, j, nil
}
