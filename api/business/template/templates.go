package templates

import (
	"strings"

	firebase_db "github.com/RemeJuan/lattr/infrastructure/firebase-db"
)

type Template struct {
	Name           string `firestore:"name"`
	OwnerId        string `firestore:"ownerId"`
	TemplateString string `firestore:"templateString"`
}

func CreateTemplate(userId string, template Template) error {
	ctx, client := firebase_db.Client()
	doc := client.Doc("templates/")

	template.OwnerId = userId
	_, err := doc.Create(ctx, template)

	return err
}

func ProcessTemplate(template string, webhookData map[string]string) string {
	for k, v := range webhookData {
		template = strings.Replace(template, "{{"+k+"}}", v, -1)
	}

	return template
}
