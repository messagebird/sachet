package messagebird

import (
	"log"
	"os"
	"fmt"

	"github.com/messagebird/go-rest-api"
	"github.com/messagebird/sachet"
)

type MessageBirdConfig struct {
	AccessKey string `yaml:"access_key"`
	Gateway   int    `yaml:"gateway"`
	Debug     bool   `yaml:"debug"`
	Language  string `yaml:"language"`
	Voice     string `yaml:"voice"`
	Repeat    int    `yaml:"repeat"`
}

type MessageBird struct {
	client *messagebird.Client
	messageParams messagebird.MessageParams
	voiceMessageParams messagebird.VoiceMessageParams
}

func NewMessageBird(config MessageBirdConfig) *MessageBird {
	client := messagebird.New(config.AccessKey)
	if config.Debug {
		client.DebugLog = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)
	}
	return &MessageBird{
		client: client,
		messageParams: messagebird.MessageParams{
			Gateway: config.Gateway,
		},
		voiceMessageParams: messagebird.VoiceMessageParams{
			Language: config.Language,
			Voice: config.Voice,
			Repeat: config.Repeat,
		},
	}
}

func (mb *MessageBird) Send(message sachet.Message) error {
	var err error=nil
	switch message.Type {
	case "","text":
		_, err = mb.client.NewMessage(message.From, message.To, message.Text, &mb.messageParams)
	case "voice":
		_, err = mb.client.NewVoiceMessage(message.To, message.Text, &mb.voiceMessageParams)
	default:
		return fmt.Errorf("unknown message type %s", message.Type)
	}
	return err
}
