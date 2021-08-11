package endpoints

import (
	"log"

	"github.com/RemeJuan/lattr/core/controller"
	"github.com/gin-gonic/gin"
)

func Register(con *controller.DBController) {
	r := gin.Default()
	g := r.Group("/")
	{
		g.POST("/create", con.Tweet)
		g.POST("/webhook", con.WebHook)
	}
	log.Fatalln(r.Run())
}

func GetParam(c *gin.Context, paramName string) string {
	return c.Params.ByName(paramName)
}
