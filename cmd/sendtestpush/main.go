package main

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	app, _ := firebase.NewApp(ctx, nil, option.WithCredentialsFile("/app/secrets/fcmServiceAccountKey.json"))
	client, _ := app.Messaging(ctx)

	message := "Some message"

	msg := &messaging.Message{
		Topic: "reminders",
		Data: map[string]string{
			"body": message,
		},
	}

	id, _ := client.Send(ctx, msg)
	fmt.Println("Message ID:", id)
}
