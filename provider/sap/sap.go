package sap

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/messagebird/sachet"
)

// Config is the configuration struct for Sap provider
type Config struct {
	URL      string `yaml:"url"`
	AuthHash string `yaml:"auth_hash"`
}

// Sap contains the necessary values for the Sap provider
type Sap struct {
	Config
	HTTPClient *http.Client // The HTTP client to send requests on
}

// NewSap creates and returns a new Sap struct
func NewSap(config Config) *Sap {
	if config.URL == "" {
		config.URL = "https://sms-pp.sapmobileservices.com/cmn/xxxxxxxxxx/xxxxxxxxxxx.sms"
	}
	return &Sap{
		config,
		&http.Client{Timeout: time.Second * 20},
	}
}

// Send sends SMS to user registered in configuration
func (c *Sap) Send(message sachet.Message) error {
	// No \n in Text tolerated
	msg := strings.ReplaceAll(message.Text, "\n", " - ")
	content := fmt.Sprintf("Version=2.0\nSubject=Alert\n[MSISDN]\nList=%s\n[MESSAGE]\nText=%s\n[SETUP]\nSplitText=yes\n[END]",
		strings.Join(message.To, ","), msg)

	request, err := http.NewRequest("POST", c.URL, strings.NewReader(content))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Basic "+c.AuthHash)

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusOK && err == nil {
		return nil
	}

	return fmt.Errorf("Failed sending sms. statusCode: %d", response.StatusCode)
}
