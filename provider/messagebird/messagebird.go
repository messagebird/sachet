package messagebird

import (
	"github.com/messagebird/go-rest-api"
	"github.com/messagebird/sachet"
)

type MessageBirdConfig struct {
	AccessKey string `yaml:"access_key"`
	Gateway   int    `yaml:"gateway"`
}

type MessageBird struct {
	client *messagebird.Client
	params messagebird.MessageParams
}

func NewMessageBird(config MessageBirdConfig) *MessageBird {
	return &MessageBird{
		client: messagebird.New(config.AccessKey),
		params: messagebird.MessageParams{
			Gateway: config.Gateway,
		},
	}
}

func (mb *MessageBird) Send(message sachet.Message) error {
	_, err := mb.client.NewMessage(message.From, message.To, message.Text, &mb.params)
	return err
}
