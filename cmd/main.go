package main

import (
	"log"

	"github.com/okassov/pet-auth/config"
	"github.com/okassov/pet-auth/internal/app"
)

var TracerName = "fib"

func main() {

	// Init Config
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	app.Run(config)
}
