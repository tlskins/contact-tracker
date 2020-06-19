package main

import (
	"log"
	"net/http"
	"os"

	d "github.com/contact-tracker/apiService/users/deliveries/http"
)

func main() {
	port := os.Getenv("USERS_PORT")

	router, err := d.Routes()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Running users on port: ", port)
	log.Panic(http.ListenAndServe(":"+port, router))
}
