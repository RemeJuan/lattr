package endpoints

import (
	"log"

	"github.com/RemeJuan/lattr/infrastructure/web-hooks"
	"github.com/gin-gonic/gin"
)

func Register() {
	r := gin.Default()
	g := r.Group("/")
	{
		g.GET("/templates/:id", templates)
		g.POST("/webhook", web_hooks.HandleWebhook)
	}
	log.Fatalln(r.Run(":8080"))
}

func templates(c *gin.Context) {
	req := getParam(c, "id")
}

func getParam(c *gin.Context, paramName string) string {
	return c.Params.ByName(paramName)
}
