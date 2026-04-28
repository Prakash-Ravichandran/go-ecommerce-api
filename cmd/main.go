package main

import (
	"log"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
	}

	api := application{
		config: cfg,
	}

	h := api.mount()
	err := api.run(h)
	if err != nil {
		log.Printf("Server has failed to start: %s", err)
		os.Exit(1)
	}
}
