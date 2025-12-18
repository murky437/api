package firebase

import (
	"context"
	"murky_api/internal/config"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type MessageService interface {
	Send(message string) error
}

type defaultMessageService struct {
	ctx    context.Context
	client *messaging.Client
}

func NewMessageService(conf *config.Config) MessageService {
	ctx := context.TODO()

	app, _ := firebase.NewApp(ctx, nil, option.WithCredentialsFile(conf.GoogleServiceAccountKeyFilePath))
	client, _ := app.Messaging(ctx)

	return &defaultMessageService{
		ctx:    ctx,
		client: client,
	}
}

func (ms *defaultMessageService) Send(message string) error {
	msg := &messaging.Message{
		Topic: "reminders",
		Data: map[string]string{
			"body": message,
		},
	}

	_, err := ms.client.Send(ms.ctx, msg)

	return err
}
