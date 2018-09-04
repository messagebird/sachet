package aspsms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

// Config is the configuration struct for AspSms provider
type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// AspSms contains the necessary values for the AspSms provider
type AspSms struct {
	Config
}

var aspsmsHTTPClient = &http.Client{Timeout: time.Second * 20}

// NewAspSms creates and returns a new AspSms struct
func NewAspSms(config Config) *AspSms {
	return &AspSms{config}
}

type requestPayload struct {
	Username    string   `json:"UserName"`
	Password    string   `json:"Password"`
	Originator  string   `json:"Originator"`
	Recipients  []string `json:"Recipients"`
	MessageText string   `json:"MessageText"`
}

// Send sends SMS to user registered in configuration
func (c *AspSms) Send(message sachet.Message) error {
	params := requestPayload{
		Username:    c.Username,
		Password:    c.Password,
		Originator:  message.From,
		Recipients:  message.To,
		MessageText: message.Text,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", "https://json.aspsms.com/SendSimpleTextSMS", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Sachet")

	response, err := aspsmsHTTPClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusOK && err == nil {
		return nil
	}

	var body []byte
	defer response.Body.Close()
	response.Body.Read(body)

	return fmt.Errorf("SMS sending failed. HTTP status code: %d, Response body: %s", response.StatusCode, body)
}
