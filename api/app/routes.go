package app

import (
	"log"
	"os"

	"github.com/RemeJuan/lattr/controllers"
	_ "github.com/RemeJuan/lattr/docs"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title lattr API
// @version 1.0
// @description API driven Tweet scheduler written in Go

// @contact.name Reme
// @contact.url https://github.com/RemeJuan/lattr
// @contact.email dev@lattr.app

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host https://api.lattr.app
// @BasePath /api/

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @scope.tweet:create Grants write access
// @scope.tweet:read Grants read and write access to administrative information

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

	if os.Getenv("GIN_MODE") != "release" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Fatalln(r.Run())
}
