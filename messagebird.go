package sachet

import "github.com/messagebird/go-rest-api"

type MessageBirdConfig struct {
	AccessKey string `yaml:"access_key"`
}

type MessageBird struct {
	client *messagebird.Client
}

func NewMessageBird(config MessageBirdConfig) *MessageBird {
	return &MessageBird{client: messagebird.New(config.AccessKey)}
}

func (mb *MessageBird) Send(message Message) error {
	_, err := mb.client.NewMessage(message.From, message.To, message.Text, nil)
	return err
}
