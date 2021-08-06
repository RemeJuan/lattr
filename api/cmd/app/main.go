package main

import (
	"github.com/RemeJuan/lattr/internal/web-hooks"
	"log"
	"net/http"
)

func main() {
	log.Println("server started")
	http.HandleFunc("/webhook", webhook.HandleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
