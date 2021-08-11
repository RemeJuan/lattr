package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RemeJuan/lattr/business/tweet"
	"github.com/RemeJuan/lattr/infrastructure/web-hooks"
	"github.com/gin-gonic/gin"
)

func Register() {
	r := gin.Default()
	g := r.Group("/")
	{
		g.POST("/create", handleTweet)
		g.POST("/webhook", web_hooks.HandleWebhook)
	}
	log.Fatalln(r.Run())
}

func handleTweet(c *gin.Context) {
	payload := make(map[string]string)
	err := json.NewDecoder(c.Request.Body).Decode(&payload)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	t, terr := tweet.BuildTweet(payload)

	if terr != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	tweet.ScheduleTweet(t)
}

func getParam(c *gin.Context, paramName string) string {
	return c.Params.ByName(paramName)
}
