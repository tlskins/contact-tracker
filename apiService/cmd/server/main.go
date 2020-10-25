package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

	chkServer, checkInHandler, err := chk.NewServer(checkInsPort, checkInsMongoDBName, checkInsMongoHost, checkInsMongo, checkInsMongoPwd, "http://localhost:"+usersPort, "http://localhost:"+placesPort, jwtKeyPath, jwtSecretPath, rpcPwd)
	if err != nil {
		log.Panic(err)
	}
	go chkServer.Start()

	placesServer, placesHandler, err := places.NewServer(placesPort, placesMongoDBName, placesMongoHost, placesMongo, placesMongoPwd, "http://localhost:"+placesPort, jwtKeyPath, jwtSecretPath, sesAccessKey, sesAccessSecret, sesRegion, senderEmail, rpcPwd)
	if err != nil {
		log.Panic(err)
	}
	go placesServer.Start()

	usersServer, usersHandler, err := users.NewServer(usersPort, usersMongoDBName, usersMongoHost, usersMongo, usersMongoPwd, jwtKeyPath, jwtSecretPath, sesAccessKey, sesAccessSecret, sesRegion, senderEmail, rpcPwd)
	if err != nil {
		log.Panic(err)
	}
	go usersServer.Start()

	apigatewayServer := gateway.NewServer(apigatewayPort, checkInsPort, placesPort, usersPort)
	go apigatewayServer.Start()

	// CLI

	ctx := context.TODO()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("---------------------")
	fmt.Printf("\nWelcome to Contact Tracker CLI\n")
	printCommands()

	for {
		fmt.Printf("\n-> ")
		command, _ := reader.ReadString('\n')
		command = strings.Replace(command, "\n", "", -1)

		if strings.Compare("stores", command) == 0 {
			places, err := placesHandler.Usecase.GetAll(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}
			fmt.Printf("%d store(s)\n", len(places))
			for _, place := range places {
				fmt.Printf("%s - %s\n", place.Name, place.Email)
			}
		} else if strings.Compare("customers", command) == 0 {
			users, err := usersHandler.Usecase.GetAll(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}
			fmt.Printf("%d customer(s)\n", len(users))
			for _, user := range users {
				fmt.Printf("%s - %s\n", user.Name, user.Email)
			}
		} else if strings.Compare("histories", command) == 0 {
			histories, err := checkInHandler.Usecase.GetHistory(ctx, "")
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}
			fmt.Printf("%d histories(s)\n", len(histories))
			for _, history := range histories {
				start := "N/A"
				if history.In != nil {
					start = history.In.Format("Jan 2 3:04 PM")
				}
				end := "N/A"
				if history.Out != nil {
					end = history.Out.Format("Jan 2 3:04 PM")
				}
				fmt.Printf("%s - %s From: %s To: %s (%d contacts)\n", history.User.Name, history.Place.Name, start, end, len(history.Contacts))
				for _, contact := range history.Contacts {
					cStart := "N/A"
					if contact.In != nil {
						cStart = contact.In.Format("Jan 2 3:04 PM")
					}
					cEnd := "N/A"
					if contact.Out != nil {
						cEnd = contact.Out.Format("Jan 2 3:04 PM")
					}
					fmt.Printf("\t%s From: %s To: %s\n", contact.User.Name, cStart, cEnd)
				}
			}
		} else if strings.Compare("help", command) == 0 {
			printCommands()
		} else {
			fmt.Printf("Invalid command\n")
		}
	}
}

func printCommands() {
	fmt.Printf("Commands:\n")
	fmt.Printf("stores : prints all available store locations\n")
	fmt.Printf("customers : prints all prior customers for all store locations\n")
	fmt.Printf("histories : prints contact histories for all customers in all store locations\n")
	fmt.Printf("help : prints these commands again\n")
}
