package app

import (
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
	}
}
