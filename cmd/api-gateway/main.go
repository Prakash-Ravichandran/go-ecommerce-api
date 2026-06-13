package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNum string = ":8080"

func main() {
	log.Println("starting our api gateway server")

	log.Println("Started on port", portNum)
	fmt.Println("To close connection CTRL+C :-)")

	// spinning up the server
	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
