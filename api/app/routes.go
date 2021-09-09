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
		tw.POST("/create", controllers.AuthenticateMiddleware("tweet:create"), controllers.CreateTweet)
		tw.GET("/:id", controllers.AuthenticateMiddleware("tweet:read"), controllers.GetTweet)
		tw.GET("/all/:userId", controllers.AuthenticateMiddleware("tweet:read"), controllers.GetTweets)
		tw.PUT("/:id", controllers.AuthenticateMiddleware("tweet:update"), controllers.UpdateTweet)
		tw.DELETE("/:id", controllers.AuthenticateMiddleware("tweet:delete"), controllers.DeleteTweet)
	}
	r.POST("/webhook", controllers.AuthenticateMiddleware("tweet:create"), controllers.WebHook)

	tk := r.Group("/token")
	{
		tk.POST("/create", controllers.TokenCreateMiddleWare("token:create"), controllers.CreateToken)
		tk.GET("/:id", controllers.AuthenticateMiddleware("token:read"), controllers.GetToken)
		tk.GET("/list", controllers.AuthenticateMiddleware("token:read"), controllers.GetTokens)
		tk.PUT("/:id", controllers.AuthenticateMiddleware("token:update"), controllers.ResetToken)
		tk.DELETE("/:id", controllers.AuthenticateMiddleware("token:delete"), controllers.DeleteToken)
	}

	log.Fatalln(r.Run())
}
