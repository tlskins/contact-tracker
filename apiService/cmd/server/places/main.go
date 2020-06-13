package main

import (
	"log"
	"net/http"
	"os"

	transport "github.com/contact-tracker/apiService/places/transport/http"
)

func main() {
	port := os.Getenv("PLACES_PORT")

	router, err := transport.Routes()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Running places on port: ", port)
	log.Panic(http.ListenAndServe(":"+port, router))
}
