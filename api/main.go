package main

import (
	"log"
	"os"

	"github.com/RemeJuan/lattr/app"
	"github.com/RemeJuan/lattr/domain/tweets"
	"github.com/RemeJuan/lattr/utils/scheduler"
)

func main() {
	log.Println("server started")

	tweets.TweetRepo.Initialize()

	if os.Getenv("GIN_MODE") == "release" {
		scheduler.Scheduler()
	}

	// setup sentry logging
	app.Init()

	// Register all available endpoints
	app.Router()
}
