package main

import "github.com/messagebird/go-rest-api"

type MessageBird struct{}

func (*MessageBird) Send(message Message) {
	client := messagebird.New(config.Providers.MessageBird.AccessKey)
	client.NewMessage(
		message.From,
		[]string{message.To},
		message.Text,
		nil)
	// TODO some error checking?
}
