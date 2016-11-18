package main

import "github.com/messagebird/go-rest-api"

type MessageBird struct{}

func (*MessageBird) Send(message Message) (err error) {
	client := messagebird.New(config.Providers.MessageBird.AccessKey)
	_, err = client.NewMessage(
		message.From,
		message.To,
		message.Text,
		nil)
	return
}
