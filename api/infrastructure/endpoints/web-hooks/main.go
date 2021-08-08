package web_hooks

import (
	"encoding/json"
	"net/http"

	templates "github.com/RemeJuan/lattr/business/template"
	"github.com/RemeJuan/lattr/infrastructure/twitter-client"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusBadRequest)
		return
	}

	webhookData := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&webhookData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(webhookData["template"]) > 0 {
		twitter_client.CreateTweet(handleTemplate(webhookData))
	} else if len(webhookData["data"]) > 0 {
		twitter_client.CreateTweet(webhookData["data"])
	} else {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
	}
}

func handleTemplate(webhookData map[string]string) string {
	return templates.ProcessTemplate("{{Title}} via /r/{{Subreddit}} {{PostURL}}", webhookData)
}
