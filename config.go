package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ReceiverConf struct {
	Name     string
	Provider string
	To       string
	From     string
	Text     string
}

var config struct {
	Providers struct {
		MessageBird struct {
			AccessKey string
		}
		Nexmo struct {
			APIKey    string
			APISecret string
		}
	}
	Receivers []ReceiverConf
}

func init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Error reading config.yaml")
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal("Error parsing config.yaml")
	}
}
