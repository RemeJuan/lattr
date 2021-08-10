package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RemeJuan/lattr/api/business/template"
	"github.com/RemeJuan/lattr/api/infrastructure/web-hooks"
	"github.com/gin-gonic/gin"
)

func Register() {
	r := gin.Default()
	g := r.Group("/")
	{
		g.GET("/templates/:id", handleTemplates)
		g.POST("/templates", handleTemplates)
		g.POST("/webhook", web_hooks.HandleWebhook)
	}
	log.Fatalln(r.Run(":8080"))
}

func handleTemplates(c *gin.Context) {
	if c.Request.Method == "GET" {
		id := getParam(c, "id")
		fmt.Println(len(id))
		fmt.Println(id)

		c.JSON(200, templates.GetTemplate(id))
	} else if c.Request.Method == "POST" {
		webhookData := make(map[string]string)
		err := json.NewDecoder(c.Request.Body).Decode(&webhookData)

		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(webhookData["data"])
		templates.CreateTemplate("001", templates.Template{
			Name:           webhookData["name"],
			TemplateString: webhookData["data"],
		})
	} else {
		http.Error(c.Writer, "Invalid payload", http.StatusBadRequest)
	}
}

func getParam(c *gin.Context, paramName string) string {
	return c.Params.ByName(paramName)
}
