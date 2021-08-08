package endpoints

import (
	"net/http"

	"github.com/RemeJuan/lattr/infrastructure/endpoints/web-hooks"
)

func Register() {
	http.HandleFunc("/webhook", web_hooks.HandleWebhook)
}
