package kavenegar

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/messagebird/sachet"
)

// Retrieving required data from 'kavenegar' sections of config.yaml.
type Config struct {
	APIToken     string   `yaml:"api_token"`
	PhoneNumbers []string `yaml:"phone_numbers"`
}

var _ (sachet.Provider) = (*KaveNegar)(nil)

// Creating the KaveNegar to contain provider data.
type KaveNegar struct {
	Config
	HTTPClient *http.Client // The HTTP client to send requests on.
}

// KaveNegar creates and returns a new KaveNegar struct.
func NewKaveNegar(config Config) *KaveNegar {
	return &KaveNegar{
		config,
		&http.Client{Timeout: time.Second * 20},
	}
}

// Building the API and call the KaveNegar endpoint to send SMS to the configured receptor from config.yaml.
func (ns *KaveNegar) Send(message sachet.Message) error {
	url := "https://api.kavenegar.com/v1/" + ns.APIToken + "/sms/send.json"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "Sachet")
	params := request.URL.Query()
	params.Add("receptor", strings.Join(message.To, ","))
	// "params.Add("sender", message.From)" retrieves the sender number using "from" under receivers section,
	// if you leave that empty, KaveNegar will use default sender SMS number to send the message.
	params.Add("sender", message.From)
	params.Add("message", message.Text)
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
