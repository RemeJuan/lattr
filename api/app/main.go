package main

import (
	"log"

	"github.com/RemeJuan/lattr/infrastructure/endpoints"
	"github.com/RemeJuan/lattr/infrastructure/sentry"
)

func main() {
	log.Println("server started")

	// Register all available endpoints
	endpoints.Register()

	// setup sentry logging
	sentry_setup.Init()
}
