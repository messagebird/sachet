package sachet

import "gopkg.in/njern/gonexmo.v1"

type NexmoConfig struct {
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
}

type Nexmo struct {
	client *nexmo.Client
}

func NewNexmo(config NexmoConfig) (*Nexmo, error) {
	client, err := nexmo.NewClientFromAPI(config.APIKey, config.APISecret)
	if err != nil {
		return nil, err
	}

	return &Nexmo{client: client}, nil
}

func (nx *Nexmo) Send(message Message) error {
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
