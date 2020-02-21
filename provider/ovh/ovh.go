package ovh

import (
	"fmt"

	"github.com/ovh/go-ovh/ovh"
	"github.com/messagebird/sachet"
)



type Config struct {
	Endpoint             string `yaml:"endpoint"`
	ApplicationKey       string `yaml:"application_key"`
	ApplicationSecret    string `yaml:"application_secret"`
	ConsumerKey          string `yaml:"consumer_key"`

	ServiceName          string `yaml:"service_name"`
}

type Ovh struct {
	client *ovh.Client
	config *Config
}

func NewOvh(config Config) (*Ovh, error) {
	client, err := ovh.NewClient(
		config.Endpoint,
		config.ApplicationKey,
		config.ApplicationSecret,
		config.ConsumerKey,
	)
	if err != nil {
		return nil, err
	}
	return &Ovh{
		client: client,
		config: &config,
	}, nil
}

func (ovh *Ovh) Send(message sachet.Message) error {
	var err error = nil
	switch message.Type {
	case "", "text":
		//receivers := make([]string, 0)
		//for _, recipent := range message.To {
		//	receivers = append(r, recipent)
		//}

		type ovhSMS map[string]interface{}
		sms := make(ovhSMS)
		sms["message"] = message.Text
		sms["noStopClause"] = false
		sms["senderForResponse"] = true
		sms["receivers"] = message.To
		//fmt.Println(sms)
		//response := make(map[string]interface{})
		//err := client.Post("/sms/sms-pn13165-1/jobs", sms, &response)
		//if err != nil {
		//	fmt.Println(err)
		//}
		if err := ovh.client.Post("/sms/sms-pn13165-1/jobs", sms, nil); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown message type %s", message.Type)
	}
	return err
}
