package firebase_db

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

func Auth() {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_KEY_PATH"))
	_, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		_ = fmt.Errorf("error initializing app: %v", err)
		return
	}
}

func Client() (context.Context, *firestore.Client) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("FIREBASE_PROJECT_ID"))

	if err != nil {
		_ = fmt.Errorf("error creating client: %v", err)
	}

	return ctx, client
}
