package app

import (
	"fmt"
	"log"
	"os"

	"github.com/getsentry/sentry-go"
)

func Init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DNS"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	} else {
		fmt.Printf("Sentry Initialized")
	}
}
