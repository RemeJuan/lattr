package web_hooks

import (
	"encoding/json"
	"github.com/RemeJuan/lattr/infrastructure/twitter-client"
	"net/http"
	"strings"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var tweet string
	webhookData := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&webhookData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(webhookData["template"]) > 0 {
		tweet = handleTemplate(webhookData)
	} else {
		tweet = webhookData["data"]
	}

	twitter_client.CreateTweet(tweet)
}

func handleTemplate(webhookData map[string]string) string {
	//get template
	template := "{{Title}} via /r/{{Subreddit}} {{PostURL}}"

	for k, v := range webhookData {
		template = strings.Replace(template, "{{"+k+"}}", v, -1)
	}

	return template
}
