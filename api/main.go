package main

import (
	"log"
	"os"

	"github.com/RemeJuan/lattr/app"
	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/scheduler"
)

// @title lattr API
// @version 1.0
// @description API driven Tweet scheduler written in Go

// @contact.name Reme
// @contact.url https://github.com/RemeJuan/lattr
// @contact.email dev@lattr.app

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host https://api.lattr.app
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @scope.tweet:create Grants write access
// @scope.tweet:read Grants read and write access to administrative information

func main() {
	log.Println("server started")

	domain.TweetRepo.Initialize()
	domain.TokenRepo.Initialize()

	if os.Getenv("GIN_MODE") == "release" {
		scheduler.Scheduler()
	}

	// setup sentry logging
	app.Init()

	// Register all available endpoints
	app.Router()
}
