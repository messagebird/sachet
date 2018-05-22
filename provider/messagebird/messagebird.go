package messagebird

import (
	"log"
	"os"

	"github.com/messagebird/go-rest-api"
	"github.com/messagebird/sachet"
)

type MessageBirdConfig struct {
	AccessKey string `yaml:"access_key"`
	Gateway   int    `yaml:"gateway"`
	Debug     bool   `yaml:"debug"`
}

type MessageBird struct {
	client *messagebird.Client
	params messagebird.MessageParams
}

func NewMessageBird(config MessageBirdConfig) *MessageBird {
	client := messagebird.New(config.AccessKey)
	if config.Debug {
		client.DebugLog = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)
	}
	return &MessageBird{
		client: client,
		params: messagebird.MessageParams{
			Gateway: config.Gateway,
		},
	}
}

func (mb *MessageBird) Send(message sachet.Message) error {
	_, err := mb.client.NewMessage(message.From, message.To, message.Text, &mb.params)
	return err
}
