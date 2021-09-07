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
		tk.POST("/create", controllers.TokenCreateMiddleWare("token:create"), controllers.CreateToken)
		tk.GET("/:id", controllers.AuthenticateMiddleware("token:read"), controllers.GetToken)
		tk.GET("/list", controllers.AuthenticateMiddleware("token:read"), controllers.GetTokens)
		tk.PUT("/:id", controllers.AuthenticateMiddleware("token:update"), controllers.ResetToken)
		tk.DELETE("/:id", controllers.AuthenticateMiddleware("token:delete"), controllers.DeleteToken)
	}

	r.POST("/webhook", controllers.WebHook)
	log.Fatalln(r.Run())
}
