package main

import twilio "github.com/carlosdp/twiliogo"

type Twilio struct{}

func (*Twilio) Send(message Message) (err error) {
	twilioClient := twilio.NewClient(config.Providers.Twilio.AccountSID, config.Providers.Twilio.AuthToken)

	for _, recipient := range message.To {
		_, err = twilio.NewMessage(
			twilioClient,
			message.From,
			recipient,
			twilio.Body(message.Text),
		)
		if err != nil {
			return
		}
	}
	return
}
