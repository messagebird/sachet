package sipgate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

const sipgateURL = "https://api.sipgate.com/v2/sessions/sms"

var sipgateHTTPClient = &http.Client{Timeout: time.Second * 20}

// Config is the configuration struct for Sipgate provider.
type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Sipgate contains the necessary values for the Sipgate provider.
type Sipgate struct {
	Config
}

// NewSipgate creates and returns a new Sipgate struct.
func NewSipgate(config Config) *Sipgate {
	return &Sipgate{config}
}

type payload struct {
	SmsID     string `json:"smsId"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

// Send sends SMS to user registered in configuration.
func (c *Sipgate) Send(message sachet.Message) error {
	for _, recipient := range message.To {
		params := payload{
			SmsID:     message.From,
			Recipient: recipient,
			Message:   message.Text,
		}

		data, err := json.Marshal(params)
		if err != nil {
			return err
		}

		request, err := http.NewRequest("POST", sipgateURL, bytes.NewBuffer(data))
		if err != nil {
			return err
		}

		request.SetBasicAuth(c.Username, c.Password)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("User-Agent", "Sachet")

		response, err := sipgateHTTPClient.Do(request)
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusNoContent {
			return fmt.Errorf("Failed sending sms. statusCode: %d", response.StatusCode)
		}
	}

	return nil
}
