package main

import "gopkg.in/njern/gonexmo.v1"

type Nexmo struct{}

func (*Nexmo) Send(message Message) (err error) {
	nexmoClient, err := nexmo.NewClientFromAPI(config.Providers.Nexmo.APIKey, config.Providers.Nexmo.APISecret)
	if err != nil {
		return
	}

	for _, recipent := range message.To {

		msg := &nexmo.SMSMessage{
			From:  message.From,
			To:    recipent,
			Type:  nexmo.Text,
			Text:  message.Text,
			Class: nexmo.Standard,
		}
		_, err = nexmoClient.SMS.Send(msg)
	}
	return
}
