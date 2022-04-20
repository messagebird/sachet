package aspsms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

// Config is the configuration struct for AspSms provider.
type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var _ (sachet.Provider) = (*AspSms)(nil)

// AspSms contains the necessary values for the AspSms provider.
type AspSms struct {
	Config

	httpClient *http.Client
}

// NewAspSms creates and returns a new AspSms struct.
func NewAspSms(config Config) *AspSms {
	return &AspSms{
		config,
		&http.Client{Timeout: time.Second * 20},
	}
}

type requestPayload struct {
	Username    string   `json:"UserName"`
	Password    string   `json:"Password"`
	Originator  string   `json:"Originator"`
	Recipients  []string `json:"Recipients"`
	MessageText string   `json:"MessageText"`
}

const apiUrl = "https://json.aspsms.com/SendSimpleTextSMS"

// Send sends SMS to user registered in configuration.
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

	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Sachet")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("SMS sending failed. HTTP status code: %d, Response body: %s", response.StatusCode, body)
	}

	return nil
}
