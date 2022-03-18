package twilio

import (
	"github.com/carlosdp/twiliogo"

	"github.com/messagebird/sachet"
)

type Config struct {
	AccountSID string `yaml:"account_sid"`
	AuthToken  string `yaml:"auth_token"`
}

var _ (sachet.Provider) = (*Twilio)(nil)

type Twilio struct {
	client twiliogo.Client
}

func NewTwilio(config Config) *Twilio {
	return &Twilio{client: twiliogo.NewClient(config.AccountSID, config.AuthToken)}
}

func (tw *Twilio) Send(message sachet.Message) error {
	for _, recipient := range message.To {
		_, err := twiliogo.NewMessage(tw.client, message.From, recipient, twiliogo.Body(message.Text))
		if err != nil {
			return err
		}
	}

	return nil
}
