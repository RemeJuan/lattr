package templates

import (
	"strings"

	"github.com/RemeJuan/lattr/infrastructure/postgress-db"
	"github.com/jinzhu/gorm"
)

type Template struct {
	gorm.Model
	Name           string `json:"name"`
	OwnerId        string `json:"ownerId"`
	TemplateString string `json:"templateString"`
}

func CreateTemplate(userId string, template Template) {
	db := postgress_db.Connect()

	defer db.Close()

	db.AutoMigrate(&Template{})

	template.OwnerId = userId
	db.Create(&template)
}

func ProcessTemplate(template string, webhookData map[string]string) string {
	for k, v := range webhookData {
		template = strings.Replace(template, "{{"+k+"}}", v, -1)
	}

	return template
}

func GetTemplate(id string) Template {
	var template Template

	db := postgress_db.Connect()

	defer db.Close()

	db.First(&template, id)

	return template
}