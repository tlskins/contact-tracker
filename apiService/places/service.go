package places

import (
	"context"
	// "flag"
	"fmt"
	"io/ioutil"
	"log"
	// "os"

	"github.com/contact-tracker/apiService/pkg/auth"
	m "github.com/contact-tracker/apiService/pkg/mongo"
	r "github.com/contact-tracker/apiService/places/repository"
	t "github.com/contact-tracker/apiService/places/types"
	// "github.com/joho/godotenv"

	"go.uber.org/zap"
)

// PlaceService - is the top level signature of this service
type PlaceService interface {
	Get(ctx context.Context, id string) (*t.Place, error)
	GetAll(ctx context.Context) ([]*t.Place, error)
	Update(ctx context.Context, place *t.UpdatePlace) (*t.Place, error)
	Create(ctx context.Context, req *t.CreatePlace) (*t.Place, error)
	Delete(ctx context.Context, id string) error

	SignIn(ctx context.Context, req *t.SignInReq) (*t.Place, error)
	Confirm(ctx context.Context, id string) error
}

// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func Init(mongoDBName, mongoHost, mongoPlace, mongoPwd, placesHost, jwtKeyPath, jwtSecretPath, sesAccessKey, sesAccessSecret, sesRegion, senderEmail, rpcPwd string) (PlaceService, *auth.JWTService, error) {
	fmt.Println("Init Places Mongo Service...")

	// cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	// flag.Parse()
	// godotenv.Load(*cfgPath)
	// mongoDBName := os.Getenv("MONGO_DB_NAME")
	// mongoHost := os.Getenv("MONGO_HOST")
	// mongoPlace := os.Getenv("MONGO_USER")
	// mongoPwd := os.Getenv("MONGO_PWD")
	// placesHost := os.Getenv("PLACES_HOST")
	// jwtKeyPath := os.Getenv("JWT_KEY_PATH")
	// jwtSecretPath := os.Getenv("JWT_SECRET_PATH")
	// sesAccessKey := os.Getenv("AWS_SES_ACCESS_KEY")
	// sesAccessSecret := os.Getenv("AWS_SES_ACCESS_SECRET")
	// sesRegion := os.Getenv("AWS_SES_REGION")
	// senderEmail := os.Getenv("SENDER_EMAIL")
	// rpcPwd := os.Getenv("RPC_AUTH_PWD")

	// Init mongo repo
	mc, err := m.NewClient(mongoHost, mongoPlace, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}
	repo := r.NewMongoPlaceRepository(mc, mongoDBName)

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

	logger, _ := zap.NewProduction()

	usecase := &LoggerAdapter{
		Logger: logger,
		Usecase: &Usecase{
			Repository: repo,
			placesHost: placesHost,
		},
	}
	return usecase, j, nil
}
