package web_hooks

import (
	"encoding/json"
	"net/http"

	"github.com/RemeJuan/lattr/api/business/template"
	"github.com/RemeJuan/lattr/api/infrastructure/twitter-client"
	"github.com/gin-gonic/gin"
)

func HandleWebhook(c *gin.Context) {
	webhookData := make(map[string]string)
	err := json.NewDecoder(c.Request.Body).Decode(&webhookData)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(webhookData["template"]) > 0 {
		twitter_client.CreateTweet(handleTemplate(webhookData))
	} else if len(webhookData["data"]) > 0 {
		twitter_client.CreateTweet(webhookData["data"])
	} else {
		http.Error(c.Writer, "Invalid payload", http.StatusBadRequest)
	}
}

func handleTemplate(webhookData map[string]string) string {
	return templates.ProcessTemplate("{{Title}} via /r/{{Subreddit}} {{PostURL}}", webhookData)
}
