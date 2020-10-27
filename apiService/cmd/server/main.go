package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
		fromEmail           = os.Getenv("NOTIFICATIONS_FROM_EMAIL")
		emailPwd            = os.Getenv("EMAIL_PWD")
		smtpHost            = os.Getenv("SMTP_HOST")
		smtpPort            = os.Getenv("SMTP_PORT")
		rpcPwd              = os.Getenv("RPC_AUTH_PWD")
	)

	chkServer, checkInHandler, err := chk.NewServer(checkInsPort, checkInsMongoDBName, checkInsMongoHost, checkInsMongo, checkInsMongoPwd, "http://localhost:"+usersPort, "http://localhost:"+placesPort, jwtKeyPath, jwtSecretPath, rpcPwd)
	if err != nil {
		log.Panic(err)
	}
	go chkServer.Start()

	placesServer, placesHandler, err := places.NewServer(placesPort, placesMongoDBName, placesMongoHost, placesMongo, placesMongoPwd, "http://localhost:"+placesPort, jwtKeyPath, jwtSecretPath, rpcPwd)
	if err != nil {
		log.Panic(err)
	}
	go placesServer.Start()

	usersServer, usersHandler, err := users.NewServer(usersPort, usersMongoDBName, usersMongoHost, usersMongo, usersMongoPwd, jwtKeyPath, jwtSecretPath, fromEmail, emailPwd, smtpHost, smtpPort, rpcPwd)
	if err != nil {
		log.Panic(err)
	}
	go usersServer.Start()

	apigatewayServer := gateway.NewServer(apigatewayPort, checkInsPort, placesPort, usersPort)
	go apigatewayServer.Start()

	// CLI

	ctx := context.TODO()
	reader := bufio.NewReader(os.Stdin)
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Panic(err)
	}
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
				fmt.Printf(
					"%s - %s From: %s To: %s (%d contacts)\n",
					history.User.Name,
					history.Place.Name,
					history.In.In(loc).Format("Jan 2 3:04 PM"),
					history.Out.In(loc).Format("Jan 2 3:04 PM"),
					len(history.Contacts),
				)
				for _, contact := range history.Contacts {
					fmt.Printf(
						"\t%s From: %s To: %s\n",
						contact.User.Name,
						contact.In.In(loc).Format("Jan 2 3:04 PM"),
						contact.Out.In(loc).Format("Jan 2 3:04 PM"),
					)
				}
			}
		} else if strings.Compare("test", command) == 0 {
			searchingUser := true
			for searchingUser {
				fmt.Printf("Type in part of the user's name or email to search for a user or q to exit the test command:\n\n-> ")
				usrSearch, _ := reader.ReadString('\n')
				usrSearch = strings.Replace(usrSearch, "\n", "", -1)
				if strings.Compare("q", usrSearch) == 0 {
					searchingUser = false
					continue
				}
				users, err := usersHandler.Usecase.Search(ctx, usrSearch)
				if err != nil {
					log.Panic(err)
				}

				fmt.Printf("%d match(es)\nEnter number next to user name to select for a test positive case or press any key to search again\n", len(users))
				for i, user := range users {
					fmt.Printf("%d) %s - %s\n", i+1, user.Name, user.Email)
				}
				fmt.Printf("\n")
				usrSelect, _ := reader.ReadString('\n')
				usrSelect = strings.Replace(usrSelect, "\n", "", -1)
				usrIdx, err := strconv.Atoi(usrSelect)
				if err != nil || usrIdx > len(users) || usrIdx < 1 {
					continue
				}
				user := users[usrIdx-1]

				fmt.Printf("\nIs this user correct? %s - %s\nY - Yes\nAny other key - No\n\n", user.Name, user.Email)
				usrConfirm, _ := reader.ReadString('\n')
				usrConfirm = strings.Replace(usrConfirm, "\n", "", -1)
				if strings.Compare("Y", usrConfirm) == 0 || strings.Compare("y", usrConfirm) == 0 {
					fmt.Printf("Selected!\n")
					searchingUser = false
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
	fmt.Printf("test : test contact alert system by simulating a positive case and notifying all users contacts\n")
	fmt.Printf("help : prints these commands again\n")
}
