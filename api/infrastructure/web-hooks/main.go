package web_hooks

import (
	"encoding/json"
	"github.com/RemeJuan/lattr/infrastructure/twitter-client"
	"net/http"
	"strings"
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
	//get template
	template := "{{Title}} via /r/{{Subreddit}} {{PostURL}}"

	for k, v := range webhookData {
		template = strings.Replace(template, "{{"+k+"}}", v, -1)
	}

	return template
}
