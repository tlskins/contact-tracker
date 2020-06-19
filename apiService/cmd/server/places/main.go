package main

import (
	"log"
	"net/http"
	"os"

	d "github.com/contact-tracker/apiService/places/deliveries/http"
)

func main() {
	port := os.Getenv("PLACES_PORT")

	router, err := d.Routes()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Running places on port: ", port)
	log.Panic(http.ListenAndServe(":"+port, router))
}
