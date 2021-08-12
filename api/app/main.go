package main

import (
	"log"

	"github.com/RemeJuan/lattr/core/controller"
	"github.com/RemeJuan/lattr/infrastructure/endpoints"
	"github.com/RemeJuan/lattr/infrastructure/postgress-db"
	"github.com/RemeJuan/lattr/infrastructure/sentry"
	"github.com/dghubble/go-twitter/twitter"
)

func main() {
	log.Println("server started")

	db := postgress.Connect()
	db.AutoMigrate(&twitter.Tweet{})
	con := controller.Init(db)
	// Register all available endpoints
	endpoints.Register(con)

	// setup sentry logging
	sentry_setup.Init()
}
