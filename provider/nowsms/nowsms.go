package nowsms

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/messagebird/sachet"
)

// Config is the configuration struct for NowSms provider
type Config struct {
	User         string   `yaml:"username"`
	Password     string   `yaml:"password"`
	PhoneNumbers []string `yaml:"phone_numbers"`
}

// NowSms contains the necessary values for the NowSms provider
type NowSms struct {
	Config
	HTTPClient *http.Client // The HTTP client to send requests on
}

// NewNowSms creates and returns a new NowSms struct
func NewNowSms(config Config) *NowSms {
	return &NowSms{
		config,
		&http.Client{Timeout: time.Second * 20},
	}
}

// Send sends SMS to user registered in configuration
func (ns *NowSms) Send(message sachet.Message) error {
	const nowSmsURL = "http://sms-gateway:8800/send"

	request, err := http.NewRequest("POST", nowSmsURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "Sachet")

	params := request.URL.Query()
	params.Add("User", ns.User)
	params.Add("Password", ns.Password)
	params.Add("PhoneNumber", strings.Join(ns.PhoneNumbers, ","))
	params.Add("Text", message.Text)
	request.URL.RawQuery = params.Encode()

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
