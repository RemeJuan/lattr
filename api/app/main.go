package main

import (
	"github.com/RemeJuan/lattr/core/sentry"
	"github.com/RemeJuan/lattr/infrastructure/web-hooks"
	"log"
	"net/http"
)

func main() {
	log.Println("server started")
	http.HandleFunc("/webhook", web_hooks.HandleWebhook)
	sentry_setup.Init()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
