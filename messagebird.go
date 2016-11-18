package main

import (
	"github.com/messagebird/go-rest-api"
	"log"
)

type MessageBird struct{}

func (*MessageBird) Send(message Message) {
	client := messagebird.New(config.Providers.MessageBird.AccessKey)
	_, err := client.NewMessage(
		message.From,
		message.To,
		message.Text,
		nil)
	if err != nil {
		log.Println(err.Error())
	}
}
