package esendex

// Provider for https://developers.esendex.com/api-reference

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

const (
	timeout  = time.Second * 60
	endpoint = "https://api.esendex.com/v1.0/messagedispatcher"
)

type Config struct {
	User             string `yaml:"user"`
	ApiToken         string `yaml:"api_token"`
	AccountReference string `yaml:"account_reference"`
}

type Esendex struct {
	Config
	httpClient *http.Client
}

func NewEsendex(config Config) *Esendex {
	return &Esendex{
		config,
		&http.Client{Timeout: timeout},
	}
}

func (e *Esendex) Send(message sachet.Message) (err error) {
	for _, phoneNumber := range message.To {
		err = e.sendOne(message, phoneNumber)

		if err != nil {
			return fmt.Errorf("failed to make API call to Esendex: %w", err)
		}
	}

	return
}

// JSON example for one message. The payload may contain multiple messages. The to property may
// contain more than one phone number separated by comma.
//	{
//	"accountreference":"xxx",
//	"messages":[{
//		"from": "from",
//		"to":"phonenuber1,phonenuber2,phonenuber3",
//		"body":"Hogla, Bogla!"
//	}]
//	}
type requestPayload struct {
	AccountReference string           `json:"accountreference"`
	Messages         []requestMessage `json:"messages"`
}

type requestMessage struct {
	From string `json:"from"`
	To   string `json:"to"`
	Body string `json:"body"`
}

func (e *Esendex) sendOne(message sachet.Message, phoneNumber string) (err error) {
	params := requestPayload{
		AccountReference: e.AccountReference,
		Messages: []requestMessage{
			{
				From: message.From,
				To:   phoneNumber,
				Body: message.Text,
			},
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	var request *http.Request
	request, err = http.NewRequest("POST", endpoint, bytes.NewBuffer(data))

	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Sachet")
	request.SetBasicAuth(e.Config.User, e.Config.ApiToken)

	response, err := e.httpClient.Do(request)
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
