package main

import (
	"log"

	"github.com/RemeJuan/lattr/app"
	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/scheduler"
)

func main() {
	log.Println("server started")

	domain.TweetRepo.Initialize()

	scheduler.Scheduler()
	// setup sentry logging
	app.Init()

	// Register all available endpoints
	app.Router()
}
