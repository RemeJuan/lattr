package app

import (
	"log"

	"github.com/RemeJuan/lattr/controllers"
	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()
	g := r.Group("/tweets")
	{
		g.POST("/create", controllers.CreateTweet)
		g.GET("/:id", controllers.GetTweet)
		g.GET("/all/:userId", controllers.GetTweets)
		g.PUT("/:id", controllers.UpdateTweet)
		g.DELETE("/:id", controllers.DeleteTweet)
	}

	r.POST("/webhook", controllers.WebHook)
	log.Fatalln(r.Run())
}
