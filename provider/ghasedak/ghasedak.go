package ghasedak

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"github.com/messagebird/sachet"
)

// Retrieving required data from 'ghasedak' sections of config.yaml
type Config struct {
	APIToken     string   `yaml:"api_token"`
	PhoneNumbers []string `yaml:"phone_numbers"`
}

// Creating the KaveNegar to contain provider data
type Ghasedak struct {
	Config
	HTTPClient *http.Client // The HTTP client to send requests on
}

// Ghasedak creates and returns a new Ghasedak struct
func NewGhasedak(config Config) *Ghasedak {
	return &Ghasedak{
		config,
		&http.Client{Timeout: time.Second * 20},
	}
}

// Building the API and call the Ghasedak endpoint to send SMS to the configured receptor from config.yaml
func (ns *Ghasedak) Send(message sachet.Message) error {
	url := "https://api.ghasedak.me/v2/sms/send/pair"
	payload := strings.NewReader("message=" + message.Text + "&receptor=" + strings.Join(message.To, ","))
	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "Sachet")
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("apikey", ns.APIToken)
	request.Header.Add("cache-control", "no-cache")
	response, err := ns.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf(
			"SMS sending failed. HTTP status code: %d, Response body: %s",
			response.StatusCode,
			body,
		)
	}
	fmt.Println("Message sent: ", message.Text)
	return nil
}
