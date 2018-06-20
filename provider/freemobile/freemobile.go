package freemobile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

// FreeMobileConfig is the configuration struct for FreeMobile provider
type FreeMobileConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Url      string `yaml:"url"`
}

// FreeMobile contains the necessary values for the FreeMobile provider
type FreeMobile struct {
	FreeMobileConfig
}

var freemobileHTTPClient = &http.Client{Timeout: time.Second * 20}

// NewFreeMobile creates and returns a new FreeMobile struct
func NewFreeMobile(config FreeMobileConfig) *FreeMobile {
	if config.Url == "" {
		config.Url = "https://smsapi.free-mobile.fr/sendmsg"
	}
	return &FreeMobile{config}
}

type FreeMobilePayload struct {
	User    string `json:"user"`
	Pass    string `json:"pass"`
	Message string `json:"msg"`
}

// Send sends SMS to user registered in configuration
func (c *FreeMobile) Send(message sachet.Message) error {
	payload := FreeMobilePayload{
		User:    c.Username,
		Pass:    c.Password,
		Message: message.Text,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", c.Url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Sachet")

	response, err := freemobileHTTPClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusOK && err == nil {
		return nil
	}

	return fmt.Errorf("Failed sending sms. statusCode: %d", response.StatusCode)
}
