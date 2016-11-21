package main

import (
	"encoding/json"
	"log"

	"gopkg.in/njern/gonexmo.v1"
)

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
		response, err := nexmoClient.SMS.Send(msg)
		if err != nil {
			return err
		}

		js0n, _ := json.Marshal(response)
		log.Println(string(js0n))
	}
	return
}
