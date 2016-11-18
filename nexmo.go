package main

import (
	"log"

	"gopkg.in/njern/gonexmo.v1"
)

type Nexmo struct{}

func (*Nexmo) Send(message Message) {
	nexmoClient, _ := nexmo.NewClientFromAPI(config.Providers.Nexmo.APIKey, config.Providers.Nexmo.APISecret)

	for _, recipent := range message.To {

		msg := &nexmo.SMSMessage{
			From:  message.From,
			To:    recipent,
			Type:  nexmo.Text,
			Text:  message.Text,
			Class: nexmo.Standard,
		}
		_, err := nexmoClient.SMS.Send(msg)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
