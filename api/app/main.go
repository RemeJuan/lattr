package main

import (
	"log"

	"github.com/RemeJuan/lattr/core/controller"
	"github.com/RemeJuan/lattr/infrastructure/endpoints"
	"github.com/RemeJuan/lattr/infrastructure/postgress-db"
	"github.com/RemeJuan/lattr/infrastructure/sentry"
)

func main() {
	log.Println("server started")

	db := postgress.Connect()
	con := controller.Init(db)
	// Register all available endpoints
	endpoints.Register(con)

	// setup sentry logging
	sentry_setup.Init()
}
