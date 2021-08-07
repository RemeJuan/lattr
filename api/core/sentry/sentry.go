package sentry_setup

import (
	"github.com/getsentry/sentry-go"
	"log"
	"os"
)

func Init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DNS"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}
