package main

import (
	"log"
	"net/http"
	"os"

	d "github.com/contact-tracker/apiService/check-ins/deliveries/http"
)

func main() {
	port := os.Getenv("CHECK_INS_PORT")

	router, err := d.Routes()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Running check ins on port: ", port)
	log.Panic(http.ListenAndServe(":"+port, router))
}
