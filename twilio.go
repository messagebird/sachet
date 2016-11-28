package sachet

import twilio "github.com/carlosdp/twiliogo"

type TwilioConfig struct {
	AccountSID string `yaml:"account_sid"`
	AuthToken  string `yaml:"auth_token"`
}

type Twilio struct {
	client twilio.Client
}

func NewTwilio(config TwilioConfig) *Twilio {
	return &Twilio{client: twilio.NewClient(config.AccountSID, config.AuthToken)}
}

func (tw *Twilio) Send(message Message) error {
	for _, recipient := range message.To {
		_, err := twilio.NewMessage(tw.client, message.From, recipient, twilio.Body(message.Text))
		if err != nil {
			return err
		}
	}

	return nil
}
