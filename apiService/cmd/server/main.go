package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	chk "github.com/contact-tracker/apiService/check-ins/deliveries/http"
	gateway "github.com/contact-tracker/apiService/cmd/server/apigateway"
	places "github.com/contact-tracker/apiService/places/deliveries/http"
	users "github.com/contact-tracker/apiService/users/deliveries/http"
)

func main() {
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	var (
		checkInsPort        = os.Getenv("CHECK_INS_PORT")
		checkInsMongoDBName = os.Getenv("CHECK_INS_MONGO_DB_NAME")
		checkInsMongoHost   = os.Getenv("CHECK_INS_MONGO_HOST")
		checkInsMongo       = os.Getenv("CHECK_INS_MONGO_USER")
		checkInsMongoPwd    = os.Getenv("CHECK_INS_MONGO_PWD")
		placesPort          = os.Getenv("PLACES_PORT")
		placesMongoDBName   = os.Getenv("PLACES_MONGO_DB_NAME")
		placesMongoHost     = os.Getenv("PLACES_MONGO_HOST")
		placesMongo         = os.Getenv("PLACES_MONGO_USER")
		placesMongoPwd      = os.Getenv("PLACES_MONGO_PWD")
		usersPort           = os.Getenv("USERS_PORT")
		usersMongoDBName    = os.Getenv("USERS_MONGO_DB_NAME")
		usersMongoHost      = os.Getenv("USERS_MONGO_HOST")
		usersMongo          = os.Getenv("USERS_MONGO_USER")
		usersMongoPwd       = os.Getenv("USERS_MONGO_PWD")
		apigatewayPort      = os.Getenv("APIGATEWAY_PORT")
		jwtKeyPath          = os.Getenv("JWT_KEY_PATH")
		jwtSecretPath       = os.Getenv("JWT_SECRET_PATH")
		sesAccessKey        = os.Getenv("AWS_SES_ACCESS_KEY")
		sesAccessSecret     = os.Getenv("AWS_SES_ACCESS_SECRET")
		sesRegion           = os.Getenv("AWS_SES_REGION")
		senderEmail         = os.Getenv("SENDER_EMAIL")
		rpcPwd              = os.Getenv("RPC_AUTH_PWD")
	)

	chkServer, err := chk.NewServer(checkInsPort, checkInsMongoDBName, checkInsMongoHost, checkInsMongo, checkInsMongoPwd, "127.0.0.1:"+usersPort, "127.0.0.1:"+placesPort, jwtKeyPath, jwtSecretPath, rpcPwd)
	if err != nil {
		log.Panic(err)
	}
	go chkServer.Start()

	placesServer, err := places.NewServer(placesPort, placesMongoDBName, placesMongoHost, placesMongo, placesMongoPwd, "127.0.0.1:"+placesPort, jwtKeyPath, jwtSecretPath, sesAccessKey, sesAccessSecret, sesRegion, senderEmail, rpcPwd)
	if err != nil {
		log.Panic(err)
	}
	go placesServer.Start()

	usersServer, err := users.NewServer(usersPort, usersMongoDBName, usersMongoHost, usersMongo, usersMongoPwd, jwtKeyPath, jwtSecretPath, sesAccessKey, sesAccessSecret, sesRegion, senderEmail, rpcPwd)
	if err != nil {
		log.Panic(err)
	}
	go usersServer.Start()

	apigatewayServer := gateway.NewServer(apigatewayPort, checkInsPort, placesPort, usersPort)
	apigatewayServer.Start()
}
