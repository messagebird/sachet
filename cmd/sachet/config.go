package main

import (
	"io/ioutil"

	"github.com/messagebird/sachet"
	"gopkg.in/yaml.v2"
)

type ReceiverConf struct {
	Name     string
	Provider string
	To       []string
	From     string
}

var config struct {
	Providers struct {
		MessageBird sachet.MessageBirdConfig
		Nexmo       sachet.NexmoConfig
		Twilio      sachet.TwilioConfig
		Infobip     sachet.InfobipConfig
	}

	Receivers []ReceiverConf
}

// LoadConfig loads the specified YAML configuration file.
func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return nil
}
