package main

import (
	"log"

	"github.com/RemeJuan/lattr/app"
	"github.com/RemeJuan/lattr/domain"
)

func main() {
	log.Println("server started")

	domain.TweetRepo.Initialize()
	// Register all available endpoints
	app.Router()

	// setup sentry logging
	app.Init()
}
