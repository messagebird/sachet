package nexmo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// USSD represents the USSD API functions for sending
// USSD push and prompt messages.
type USSD struct {
	client *Client
}

// USSDMessage represents a single USSD message.
type USSDMessage struct {
	From                 string
	To                   string
	Text                 string
	StatusReportRequired bool   // Optional.
	ClientReference      string // Optional.
	NetworkCode          string // Optional.

	// Optional: If true, message will be a USSD prompt type,
	// otherwise it will be a push.
	Prompt bool
}

// Send the message using the specified USSD client.
func (c *USSD) Send(msg *USSDMessage) (*MessageResponse, error) {
	if len(msg.From) <= 0 {
		return nil, errors.New("Invalid From field specified")
	}

	if len(msg.To) <= 0 {
		return nil, errors.New("Invalid To field specified")
	}

	if len(msg.ClientReference) > 40 {
		return nil, errors.New("Client reference too long")
	}

	var messageResponse *MessageResponse

	values := make(url.Values)

	if len(msg.Text) <= 0 {
		return nil, errors.New("Invalid message text")
	} else {
		// TODO(inhies): UTF8 and URL encode before setting
		values.Set("text", msg.Text)
	}

	if !c.client.useOauth {
		values.Set("api_key", c.client.apiKey)
		values.Set("api_secret", c.client.apiSecret)
	}

	if msg.StatusReportRequired {
		values.Set("status_report_req", "1")
	}

	if msg.ClientReference != "" {
		values.Set("client_ref", msg.ClientReference)
	}

	if msg.NetworkCode != "" {
		values.Set("network-code", msg.NetworkCode)
	}

	var endpoint string
	if msg.Prompt {
		endpoint = "/ussd-prompt/json"
	} else {
		endpoint = "/ussd/json"
	}
	values.Set("to", msg.To)
	values.Set("from", msg.From)

	valuesReader := bytes.NewReader([]byte(values.Encode()))
	var r *http.Request
	r, _ = http.NewRequest("POST", apiRoot+endpoint, valuesReader)

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.HttpClient.Do(r)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &messageResponse)
	if err != nil {
		return nil, err
	}
	return messageResponse, nil
}
