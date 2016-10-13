package main

import (
	"io/ioutil"
	"log"

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
		MessageBird struct {
			AccessKey string `yaml:"access_key"`
		}
		Nexmo struct {
			APIKey    string `yaml:"api_key"`
			APISecret string `yaml:"api_secret"`
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
