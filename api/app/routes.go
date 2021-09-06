package app

import (
	"log"

	"github.com/RemeJuan/lattr/controllers"
	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()
	tw := r.Group("/tweets")
	{
		tw.POST("/create", controllers.CreateTweet)
		tw.GET("/:id", controllers.GetTweet)
		tw.GET("/all/:userId", controllers.GetTweets)
		tw.PUT("/:id", controllers.UpdateTweet)
		tw.DELETE("/:id", controllers.DeleteTweet)
	}

	tk := r.Group("/token")
	{
		tk.POST("/create", controllers.TokenCreateMiddleWare(controllers.CreateToken))
		tk.GET("/:id", controllers.GetToken)
		tk.GET("/list", controllers.GetToken)
		tk.PUT("/:id", controllers.ResetToken)
		tk.DELETE("/:id", controllers.DeleteToken)
	}

	r.POST("/webhook", controllers.WebHook)
	log.Fatalln(r.Run())
}
