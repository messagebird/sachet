package nexmo

import (
	"github.com/messagebird/sachet"

	nexmo "gopkg.in/njern/gonexmo.v1"
)

type Config struct {
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
}

var _ (sachet.Provider) = (*Nexmo)(nil)

type Nexmo struct {
	client *nexmo.Client
}

func NewNexmo(config Config) (*Nexmo, error) {
	client, err := nexmo.NewClientFromAPI(config.APIKey, config.APISecret)
	if err != nil {
		return nil, err
	}

	return &Nexmo{client: client}, nil
}

func (nx *Nexmo) Send(message sachet.Message) error {
	for _, recipent := range message.To {
		msg := &nexmo.SMSMessage{
			From:  message.From,
			To:    recipent,
			Type:  nexmo.Text,
			Text:  message.Text,
			Class: nexmo.Standard,
		}

		if _, err := nx.client.SMS.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
