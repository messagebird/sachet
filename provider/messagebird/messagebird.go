package messagebird

import (
	"fmt"
	"log"
	"os"

	messagebird "github.com/messagebird/go-rest-api"
	sms "github.com/messagebird/go-rest-api/sms"
	voicemessage "github.com/messagebird/go-rest-api/voicemessage"
	"github.com/messagebird/sachet"
)

type Config struct {
	AccessKey string `yaml:"access_key"`
	Gateway   int    `yaml:"gateway"`
	Debug     bool   `yaml:"debug"`
	Language  string `yaml:"language"`
	Voice     string `yaml:"voice"`
	Repeat    int    `yaml:"repeat"`
}

type MessageBird struct {
	client             *messagebird.Client
	messageParams      sms.Params
	voiceMessageParams voicemessage.Params
}

func NewMessageBird(config Config) *MessageBird {
	client := messagebird.New(config.AccessKey)
	if config.Debug {
		client.DebugLog = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)
	}
	return &MessageBird{
		client: client,
		messageParams: sms.Params{
			Gateway: config.Gateway,
		},
		voiceMessageParams: voicemessage.Params{
			Language: config.Language,
			Voice:    config.Voice,
			Repeat:   config.Repeat,
		},
	}
}

func (mb *MessageBird) Send(message sachet.Message) error {
	var err error = nil
	switch message.Type {
	case "", "text":
		_, err = sms.Create(mb.client, message.From, message.To, message.Text, &mb.messageParams)
	case "voice":
		_, err = voicemessage.Create(mb.client, message.To, message.Text, &mb.voiceMessageParams)
	default:
		return fmt.Errorf("unknown message type %s", message.Type)
	}
	return err
}
