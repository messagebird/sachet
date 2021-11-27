package sms77

import (
	"fmt"
	"strings"

	"github.com/messagebird/sachet"
	"github.com/sms77io/go-client/sms77api"
)

// Config is the configuration struct for Sms77 provider
type Config struct {
	ApiKey string `yaml:"api_key"`
	Debug  bool   `yaml:"debug"`
}

// Sms77 contains the necessary values for the Sms77 provider
type Sms77 struct {
	client *sms77api.Sms77API
	config Config
}

// NewSms77 creates and returns a new Sms77 struct
func NewSms77(config Config) *Sms77 {
	client := sms77api.New(sms77api.Options{
		ApiKey:   config.ApiKey,
		Debug:    config.Debug,
		SentWith: "Sachet",
	})

	return &Sms77{
		client,
		config,
	}
}

// Send sends SMS to user registered in configuration
func (s77 *Sms77) Send(message sachet.Message) error {
	var err error = nil
	switch message.Type {
	case "", "text":
		_, err = s77.client.Sms.Json(sms77api.SmsBaseParams{
			From: message.From,
			Text: message.Text,
			To:   strings.Join(message.To, ","),
		})
	case "voice":
		for _, recipient := range message.To {
			_, err = s77.client.Voice.Json(sms77api.VoiceParams{
				From: message.From,
				Text: message.Text,
				To:   recipient,
			})
		}
	default:
		return fmt.Errorf("unknown message type %s", message.Type)
	}
	return err
}
