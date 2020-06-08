package main

import (
	"log"
	"net/http"
	"os"

	transport "github.com/contact-tracker/apiService/users/transport/http"
)

func main() {
	port := os.Getenv("PORT")

	router, err := transport.Routes()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Running on port: ", port)
	log.Panic(http.ListenAndServe(":"+port, router))
}
