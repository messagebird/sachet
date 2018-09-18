package sipgate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

var sipgateHTTPClient = &http.Client{Timeout: time.Second * 20}

// Config is the configuration struct for Sipgate provider
type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
}

// Sipgate contains the necessary values for the Sipgate provider
type Sipgate struct {
	Config
}

// NewSipgate creates and returns a new Sipgate struct
func NewSipgate(config Config) *Sipgate {
	if config.URL == "" {
		config.URL = "https://api.sipgate.com/v2/sessions/sms"
	}
	return &Sipgate{config}
}

type payload struct {
	SmsID     string `json:"smsId"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

// Send sends SMS to user registered in configuration
func (c *Sipgate) Send(message sachet.Message) error {
	for _, recipient := range message.To {
		var body bytes.Buffer

		params := payload{
			SmsID:     message.From,
			Recipient: recipient,
			Message:   message.Text,
		}

		enc := json.NewEncoder(&body)
		if err := enc.Encode(params); err != nil {
			return err
		}

		request, err := http.NewRequest("POST", c.URL, &body)
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
