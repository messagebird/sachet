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

// LoadConfig parses the given YAML file into a Config.
func LoadConfig(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading configuration file")
	}

	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatal("Error parsing configuration file")
	}
}
