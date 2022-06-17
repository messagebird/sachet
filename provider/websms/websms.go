package websms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

const websmsBaseURL = "https://api.websms.com/rest/"

var httpClient = &http.Client{Timeout: 20 * time.Second}

// Config is the configuration struct for websms provider.
type Config struct {
	Token string `yaml:"api_token"`
}

var _ (sachet.Provider) = (*WebSms)(nil)

// Websms contains the necessary values for the WebSms provider.
type WebSms struct {
	Config
}

// NewWebSms creates and returns a new Websms struct.
func NewWebSms(config Config) *WebSms {
	return &WebSms{config}
}

type TextSmsSendRequest struct {
	SmsID			string   `json:"clientMessageId"`
	RecipientList	[]string `json:"recipientAddressList"`
	Message  		string   `json:"messageContent"`
}

// Send sends SMS to user registered in configuration.
func (c *WebSms) Send(message sachet.Message) error {
	
	// Prepare SMS request payload
	smsReq := TextSmsSendRequest{
		RecipientList: message.To,
		Message:   message.Text,
	}

	jsonSmsReq, err := json.Marshal(smsReq)
	if err != nil {
		return err
	}

	// Create HTTP POST request
	request, err := http.NewRequest("POST", websmsBaseURL + "smsmessaging/text", bytes.NewBuffer(jsonSmsReq))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Sachet")
	request.Header.Set("Authorization", "Bearer " + c.Token)

	// Send HTTP request
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send sms. statusCode: %d", response.StatusCode)
	}
	
	// Messages successfully sent
	return nil
}
